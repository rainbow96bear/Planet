package storage

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// localFileStorage는 로컬 디스크에 파일을 저장하는 FileStorage 구현체다.
// 로컬 개발/테스트 환경 전용이다. Cloud Run 등 파일시스템이 휘발성인
// 플랫폼에 배포할 때는 이 구현체를 쓰면 안 되며, 별도 구현체(예: Supabase
// Storage, Cloudflare R2)로 교체해야 한다.
type localFileStorage struct {
	baseDir string // 파일이 저장될 로컬 디렉토리, 예: "./uploads"
	baseURL string // 저장된 파일에 접근할 때 쓸 URL prefix, 예: "http://localhost:8080/uploads"
}

// NewLocalFileStorage는 로컬 디스크 기반 FileStorage를 생성한다.
func NewLocalFileStorage(baseDir, baseURL string) FileStorage {
	return &localFileStorage{
		baseDir: baseDir,
		baseURL: strings.TrimSuffix(baseURL, "/"),
	}
}

func (s *localFileStorage) Upload(input UploadInput) (string, error) {
	if input.Reader == nil {
		return "", errors.New("업로드할 파일이 없습니다")
	}

	fullPath, err := s.resolvePath(input.Key)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return "", fmt.Errorf("디렉토리 생성 실패: %w", err)
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("파일 생성 실패: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, input.Reader); err != nil {
		return "", fmt.Errorf("파일 쓰기 실패: %w", err)
	}

	url := fmt.Sprintf("%s/%s", s.baseURL, filepath.ToSlash(input.Key))
	return url, nil
}

func (s *localFileStorage) Delete(url string) error {
	key, err := s.extractKey(url)
	if err != nil {
		return err
	}

	fullPath, err := s.resolvePath(key)
	if err != nil {
		return err
	}

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			// 이미 없는 파일이면 조용히 넘어감 (베스트 에포트 삭제 정책과 일관됨)
			slog.Warn("delete: file already absent", "path", fullPath)
			return nil
		}
		return fmt.Errorf("파일 삭제 실패: %w", err)
	}

	return nil
}

// resolvePath는 Key를 실제 파일 경로로 변환하면서, 그 경로가 baseDir을
// 벗어나지 않는지 검증한다. Key에 "../" 같은 경로 순회 문자열이 섞여
// 들어와도(예: 조작된 업로드 파일명) baseDir 바깥을 건드릴 수 없도록
// 막는 심층 방어 지점이다.
func (s *localFileStorage) resolvePath(key string) (string, error) {
	absBase, err := filepath.Abs(s.baseDir)
	if err != nil {
		return "", fmt.Errorf("기준 경로 확인 실패: %w", err)
	}

	fullPath := filepath.Join(absBase, filepath.FromSlash(key))

	// baseDir 경계를 명확히 하기 위해 구분자를 붙여서 비교한다.
	// (예: baseDir이 "/data/uploads"일 때 "/data/uploads-evil"과 헷갈리지 않도록)
	if fullPath != absBase && !strings.HasPrefix(fullPath, absBase+string(filepath.Separator)) {
		return "", fmt.Errorf("허용되지 않은 경로입니다: %s", key)
	}

	return fullPath, nil
}

// extractKey는 baseURL을 기준으로 URL에서 저장 시 사용한 Key를 복원한다.
func (s *localFileStorage) extractKey(url string) (string, error) {
	prefix := s.baseURL + "/"
	if !strings.HasPrefix(url, prefix) {
		return "", fmt.Errorf("알 수 없는 URL 형식입니다: %s", url)
	}
	return strings.TrimPrefix(url, prefix), nil
}

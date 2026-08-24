package storage

import (
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"time"
)

// UploadInput은 파일 저장에 필요한 데이터를 담는다.
type UploadInput struct {
	Key         string
	Reader      io.Reader
	ContentType string
	Size        int64
}

// FileStorage는 파일 저장소에 대한 추상화다.
// 로컬 디스크, Supabase Storage, Cloudflare R2, GCS 등
// 어떤 구현체로 교체되든 이 인터페이스만 만족하면 서비스 레이어 코드는 변경되지 않는다.
type FileStorage interface {
	// Upload는 파일을 저장하고 접근 가능한 URL을 반환한다.
	Upload(input UploadInput) (string, error)
	// Delete는 저장된 파일을 삭제한다. url은 Upload가 반환한 값이다.
	Delete(url string) error
}

var unsafeFilenameChars = regexp.MustCompile(`[^a-zA-Z0-9._-]`)

// sanitizeFilename은 경로 순회 문자열(../ 등)과 특수문자를 제거해
// Key 생성에 안전하게 쓸 수 있는 파일명만 남긴다.
func sanitizeFilename(filename string) string {
	base := filepath.Base(filename) // 디렉토리 구분자를 제거 (예: "../../etc/passwd" → "passwd")
	safe := unsafeFilenameChars.ReplaceAllString(base, "_")
	if safe == "" || safe == "." || safe == ".." {
		safe = "file"
	}
	return safe
}

// ProfileImageKey는 유저 프로필 이미지의 저장 경로(Key)를 만든다.
// 형식: profile-images/{userID}/{unix nano timestamp}_{filename}
// 타임스탬프를 붙이는 이유는 같은 파일명으로 재업로드해도 캐시가 깨지도록 하기 위함이다.
func ProfileImageKey(userID, filename string) string {
	safeFilename := sanitizeFilename(filename)
	return fmt.Sprintf("profile-images/%s/%d_%s", userID, time.Now().UnixNano(), safeFilename)
}

// 나중에 피드/게시물 이미지 등이 추가되면 이런 식으로 Key 생성 함수만 늘어난다.
// FileStorage 인터페이스와 구현체는 건드릴 필요가 없다.
//
// func FeedImageKey(feedID, filename string) string {
// 	safeFilename := sanitizeFilename(filename)
// 	return fmt.Sprintf("feed-images/%s/%d_%s", feedID, time.Now().UnixNano(), safeFilename)
// }

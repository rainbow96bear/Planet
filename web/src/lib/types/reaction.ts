export type ReactionType = 'like' | 'cheer'

export interface ToggleReactionBody {
    type: ReactionType
}
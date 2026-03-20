import { create } from 'zustand'

interface CommentEditStore {
  editingId: number | null
  startEdit: (id: number) => void
  stopEdit: () => void
}

export const useCommentEditStore = create<CommentEditStore>((set) => ({
  editingId: null,
  startEdit: (id) => set({ editingId: id }),
  stopEdit: () => set({ editingId: null }),
}))

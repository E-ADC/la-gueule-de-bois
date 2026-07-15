/** Icônes SVG compactes pour les votes (remplacent le texte verbeux "+1/-1"). */

export function ThumbUpIcon({ size = 16 }: { size?: number }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={2} strokeLinecap="round" strokeLinejoin="round">
      <path d="M7 11v9H3v-9h4Zm0 0 4-8a2 2 0 0 1 2 2v5h5a2 2 0 0 1 2 2l-1.5 7a2 2 0 0 1-2 1.5H7" />
    </svg>
  )
}

export function ThumbDownIcon({ size = 16 }: { size?: number }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={2} strokeLinecap="round" strokeLinejoin="round">
      <path d="M17 13V4h4v9h-4Zm0 0-4 8a2 2 0 0 1-2-2v-5H6a2 2 0 0 1-2-2l1.5-7A2 2 0 0 1 7.5 3.5H17" />
    </svg>
  )
}

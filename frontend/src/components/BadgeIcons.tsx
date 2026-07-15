import type { ReactNode } from 'react'

/** Icônes SVG pour les 4 badges du système de scoring. */

interface BadgeIconProps {
  size?: number
  title?: string
}

/** Première Cuite — Chope de bière */
export function IconPremiereCuite({ size = 32, title }: BadgeIconProps) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="#b5651d"
      strokeWidth={1.5}
      strokeLinecap="round"
      strokeLinejoin="round"
      {...(title && { title })}
    >
      {/* Corps de la chope */}
      <rect x="4" y="4" width="10" height="14" rx="1" />
      {/* Anse */}
      <path d="M 14 6 Q 18 6 18 11 Q 18 16 14 16" />
      {/* Mousse */}
      <ellipse cx="9" cy="4" rx="5" ry="2" fill="#8b4513" />
    </svg>
  )
}

/** Habitué du Bar — Tabouret + Chope */
export function IconHabitueDuBar({ size = 32, title }: BadgeIconProps) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="#d2691e"
      strokeWidth={1.5}
      strokeLinecap="round"
      strokeLinejoin="round"
      {...(title && { title })}
    >
      {/* Tabouret — siège circulaire */}
      <ellipse cx="12" cy="10" rx="4" ry="2" />
      {/* Pied du tabouret */}
      <line x1="12" y1="12" x2="12" y2="18" />
      {/* Base du tabouret */}
      <circle cx="12" cy="20" r="3" fill="none" stroke="#d2691e" />
      {/* Chope superposée petite en haut à droite */}
      <rect x="16" y="2" width="5" height="6" rx="0.5" />
      <path d="M 21 3 Q 23 3 23 5 Q 23 7 21 7" strokeWidth={1} />
    </svg>
  )
}

/** Légende de la Soirée — Étoile */
export function IconLegendeDeLaSoiree({ size = 32, title }: BadgeIconProps) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="#4a2c17"
      strokeWidth={1.5}
      strokeLinecap="round"
      strokeLinejoin="round"
      {...(title && { title })}
    >
      {/* Étoile à 5 branches */}
      <polygon points="12,2 15,10 23,10 17,15 19,23 12,18 5,23 7,15 1,10 9,10" />
    </svg>
  )
}

/** Roi de la Gueule de Bois — Couronne */
export function IconRoiDeLaGueuleDeBois({ size = 32, title }: BadgeIconProps) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="#b5651d"
      strokeWidth={1.5}
      strokeLinecap="round"
      strokeLinejoin="round"
      {...(title && { title })}
    >
      {/* Bande de base de la couronne */}
      <path d="M 2 16 L 4 8 L 8 12 L 12 4 L 16 12 L 20 8 L 22 16 Z" fill="none" />
      {/* Ligne de base */}
      <line x1="2" y1="16" x2="22" y2="16" />
      {/* Points de la couronne (petits cercles) */}
      <circle cx="8" cy="12" r="1.5" fill="#b5651d" />
      <circle cx="12" cy="4" r="1.5" fill="#b5651d" />
      <circle cx="16" cy="12" r="1.5" fill="#b5651d" />
    </svg>
  )
}

/**
 * Fonction utilitaire qui retourne le composant d'icône correspondant au code du badge.
 * Retourne `null` si le code n'est pas reconnu.
 */
export function badgeIconByCode(code: string): ReactNode {
  switch (code) {
    case 'premiere-cuite':
      return <IconPremiereCuite />
    case 'habitue-du-bar':
      return <IconHabitueDuBar />
    case 'legende-de-la-soiree':
      return <IconLegendeDeLaSoiree />
    case 'roi-de-la-gueule-de-bois':
      return <IconRoiDeLaGueuleDeBois />
    default:
      return null
  }
}

export function calcAge(birthday) {
  if (!birthday) return null
  const birth = new Date(birthday)
  const today = new Date()
  let age = today.getFullYear() - birth.getFullYear()
  const monthDiff = today.getMonth() - birth.getMonth()
  if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birth.getDate())) {
    age--
  }
  return age
}

export function ageLabel(birthday) {
  const age = calcAge(birthday)
  if (age === null) return null
  const last = age % 10
  const last2 = age % 100
  if (last2 >= 11 && last2 <= 19) return `${age} лет`
  if (last === 1) return `${age} год`
  if (last >= 2 && last <= 4) return `${age} года`
  return `${age} лет`
}
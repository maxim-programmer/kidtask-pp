import { ref, computed } from 'vue'

export function getGreeting(): string {
  const now = ref(new Date())
  const tod = computed(() => now.value.toTimeString().slice(0, 8)) // HH:MM:SS
  const hour = computed((): number => {
    const hour_num = tod.value.split(':')
    return parseInt(hour_num[0] ?? '0', 10)
  })
  if (hour.value >= 6 && hour.value < 10)
    return "Доброе утро"
  if (hour.value >= 10 && hour.value < 17)
    return "Добрый день"
  if (hour.value >= 17 && hour.value < 23)
    return "Добрый вечер"
  return "Доброй ночи"
}

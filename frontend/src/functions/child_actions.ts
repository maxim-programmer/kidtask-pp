import { ref, readonly } from 'vue'

const child_user = ref<{
  username: string,
  password_hash: string,
  avatar_style: string,
  avatar_path: string,
  age_group: string,
  birthday: Date,
  parent_id: number
} | null>(null)

export function getChildUserData(uid: number) {
  console.log(uid)
  // fetch child from db by uid
  return readonly(child_user)
}

export const test_child_user1 = {
  username: 'Иван',
  password_hash: 'dj232jheiw8f',
  avatar_style: 'default-avatar',
  avatar_path: '/parent/profile.png',
  age_group: 'ch1114',
  parent_id: 1
}

export const test_child_user2 = {
  username: 'Мария',
  password_hash: 'dj232jheiw8f',
  avatar_style: 'default-avatar',
  avatar_path: '/parent/profile.png',
  age_group: 'ch710',
  parent_id: 1
}

export function check_login_credentials(): boolean {
  return true;
}

export function add_wish(title: string, description: string): void {
  console.log('add_wish', title, description)
  // TODO: отправить на бэкенд
}

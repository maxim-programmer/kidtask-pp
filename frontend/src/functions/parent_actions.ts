import { ref, readonly } from 'vue'

const parent_user = ref<{
  id: number,
  username: string,
  password_hash: string,
  avatar_style: string,
  avatar_path: string
} | null>(null)

export function getUserData(uid: number) {
  console.log(uid)
  // fetch parent from db by uid
  return readonly(parent_user)
}

export const test_user = {
  id: 1,
  username: 'Анна',
  password_hash: 'dj232jheiw8f',
  avatar_style: 'default-avatar',
  avatar_path: '/parent/profile.png'
}

export function check_login_credentials() {
  return true
}
export function check_reg_credentials() {
  return true
}

export function add_child() {
  return true
}

export function edit_child() {
  return true
}

export function remove_child() {
  return true
}

export function add_task() {
  return true
}

export function edit_task() {
  return true
}

export function mark_task_done() {
  return true
}

export function mark_task_fix() {
  return true
}

export function remove_task() {
  return true
}

export function edit_wish() {
  return true
}

export function remove_wish() {
  return true
}


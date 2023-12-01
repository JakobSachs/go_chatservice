<script setup lang="ts">
export interface ChatGroupLabelProps {
  title: string;
  last_msg: string;
  is_selected: boolean;
}

const props = withDefaults(defineProps<ChatGroupLabelProps>(), {
  title: "You shouldnt be seeing this",
  last_msg: "Hello from the dev",
  is_selected: false,
});

function restrictLength(str: string, len: number): string {
  if (str.length > len) {
    return str.slice(0, len) + "...";
  }
  return str;
}
</script>

<template>
  <div
    class="my-1 p-2 shadow rounded flex flex-row items-center hover:bg-gray-100"
  >
    <div>
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        fill="currentColor"
        class="w-6 h-6"
      >
        <path
          fill-rule="evenodd"
          d="M18.685 19.097A9.723 9.723 0 0021.75 12c0-5.385-4.365-9.75-9.75-9.75S2.25 6.615 2.25 12a9.723 9.723 0 003.065 7.097A9.716 9.716 0 0012 21.75a9.716 9.716 0 006.685-2.653zm-12.54-1.285A7.486 7.486 0 0112 15a7.486 7.486 0 015.855 2.812A8.224 8.224 0 0112 20.25a8.224 8.224 0 01-5.855-2.438zM15.75 9a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z"
          clip-rule="evenodd"
        />
      </svg>
    </div>
    <div class="flex flex-col items-baseline">
      <span class="mx-2 text-base font-bold">{{ props.title }}:</span>
      <p class="text-xs mx-2 font-light text-gray-600 dark:text-gray-400">
        {{ restrictLength(props.last_msg, 19) }}
      </p>
    </div>
  </div>
</template>

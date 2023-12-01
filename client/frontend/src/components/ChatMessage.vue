<script setup lang="ts">
export interface ChatMessageProps {
    from: string;
    text: string;
    timestamp: string;
    is_from_user: boolean;
}

const props = withDefaults(defineProps<ChatMessageProps>(), {
    from: "You shouldnt be seeing this",
    text: "Hello from the dev",
    timestamp: "2021-08-01T00:00:00Z",
    is_from_user: false,
});
</script>
<template>
    <div class="bg-gray-100 bg-opacity-20 m-2 p-4 w-full max-w-xl rounded-md flex-col items-start shadow inline-flex"
        :class="{
            'justify-end': props.is_from_user,
            'ml-auto': props.is_from_user,
            'bg-blue-700': props.is_from_user,
            'text-white': props.is_from_user,
        }">
        <p class="text-gray-700 font-medium text-sm" :class="{
            'text-gray-100': props.is_from_user,
        }">{{ from }}</p>
        <p class="relative text-black font-normal text-sm" :class="{
            'text-gray-100': props.is_from_user,
        }" v-html="text"></p>
        <p class="ml-auto text-gray-500 font-light text-xs">{{ timestamp }}</p>
    </div>
</template>

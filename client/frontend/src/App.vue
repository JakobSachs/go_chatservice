<script setup lang="ts">
import { forStatement } from "@babel/types";
import { ref } from "vue";
import ChatGroupLabel, {
    ChatGroupLabelProps,
} from "./components/ChatGroupLabel.vue";

import ChatMessage, { ChatMessageProps } from "./components/ChatMessage.vue";

const groups = ref(
    Array.from<ChatGroupLabelProps>([
        {
            title: "General",
            last_msg: "Hello from the dev ya filthy rat🐀",
            is_selected: false,
        },
        {
            title: "Your mom",
            last_msg: "ohh yeah def me baby UHHHHHHHHHHH",
            is_selected: false,
        },
    ]),
);
const selected_group = ref(0);

const messages = ref(
    Array.from<ChatMessageProps>([
        {
            from: "You",
            text: "Hello from the dev ya filthy rat🐀",
            timestamp: "1 minute ago",
            is_from_user: true,
        },
        {
            from: "Your mom",
            text: "ohh yeah def me baby UHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHhH",
            timestamp: "now",
            is_from_user: false,
        },
    ]),
);

console.log(groups);
</script>

<template>
    <div class="app bg-transparent flex flex-row min-h-screen dark:text-gray-200">
        <aside class="sidebar w-64 md:shadow transform">
            <div class="sidebar-header flex items-center justify-center py-4">
                <div class="inline-flex">
                    <a href="#" class="inline-flex flex-row items-center">
                        <svg class="h-8 w-8 mr-3 dark:fill-white" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"
                            fill="currentColor">
                            <path fill-rule="evenodd"
                                d="M4.848 2.771A49.144 49.144 0 0112 2.25c2.43 0 4.817.178 7.152.52 1.978.292 3.348 2.024 3.348 3.97v6.02c0 1.946-1.37 3.678-3.348 3.97-1.94.284-3.916.455-5.922.505a.39.39 0 00-.266.112L8.78 21.53A.75.75 0 017.5 21v-3.955a48.842 48.842 0 01-2.652-.316c-1.978-.29-3.348-2.024-3.348-3.97V6.741c0-1.946 1.37-3.68 3.348-3.97z"
                                clip-rule="evenodd" />
                        </svg>
                        <span class="leading-10 text-gray-900 dark:text-gray-100 text-l font-bold ml-1 uppercase">Wails
                            Chatclient</span>
                    </a>
                </div>
            </div>
            <div class="sidebar-content px-4 py-6">
                <!--Group list -->
                <ul class="flex flex-col w-full">
                    <!--Group header -->
                    <li class="flex flex-row items-center h-10 px-3 rounded-lg text-gray-700 dark:text-gray-200">
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"
                            class="w-6 h-6 fill-gray-700 dark:fill-gray-200">
                            <path fill-rule="evenodd"
                                d="M8.25 6.75a3.75 3.75 0 117.5 0 3.75 3.75 0 01-7.5 0zM15.75 9.75a3 3 0 116 0 3 3 0 01-6 0zM2.25 9.75a3 3 0 116 0 3 3 0 01-6 0zM6.31 15.117A6.745 6.745 0 0112 12a6.745 6.745 0 016.709 7.498.75.75 0 01-.372.568A12.696 12.696 0 0112 21.75c-2.305 0-4.47-.612-6.337-1.684a.75.75 0 01-.372-.568 6.787 6.787 0 011.019-4.38z"
                                clip-rule="evenodd" />
                            <path
                                d="M5.082 14.254a8.287 8.287 0 00-1.308 5.135 9.687 9.687 0 01-1.764-.44l-.115-.04a.563.563 0 01-.373-.487l-.01-.121a3.75 3.75 0 013.57-4.047zM20.226 19.389a8.287 8.287 0 00-1.308-5.135 3.75 3.75 0 013.57 4.047l-.01.121a.563.563 0 01-.373.486l-.115.04c-.567.2-1.156.349-1.764.441z" />
                        </svg>
                        <span class="mx-3">Groups</span>
                    </li>
                    <!-- Generate groups for each group -->
                    <ChatGroupLabel v-for="(group, i) in groups" :title="group.title" :last_msg="group.last_msg"
                        :is_selected="i == selected_group" />
                </ul>
            </div>
        </aside>
        <main class="main flex flex-col flex-grow -ml-64 md:ml-0 transition-all duration-150 ease-in">
            <div class="main-content flex flex-col flex-grow p-4">
                <h1 class="font-bold text-2xl text-gray-700 dark:text-gray-300">
                    {{ groups[selected_group].title }}
                </h1>
                <div class="flex flex-col flex-grow border-white border border-opacity-0 rounded mt-4">
                    <div class="flex flex-col flex-grow w-full max-w-screen overflow-hidden">
                        <ChatMessage v-for="msg in messages" :from="msg.from" :text="msg.text" :timestamp="msg.timestamp"
                            :is_from_user="msg.is_from_user" />
                    </div>
                </div>
            </div>
        </main>
    </div>
</template>

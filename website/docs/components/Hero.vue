<template>
<section class="bg-white dark:bg-[#1b1b1f]">
    <div class="py-8 px-4 mx-auto max-w-screen-xl text-center lg:py-16 lg:px-12">
        <a href="#" class="no-underline-important inline-flex justify-between items-center py-1 px-1 pr-4 mb-7 text-sm text-gray-700 bg-gray-100 rounded-full dark:bg-gray-800 dark:text-white hover:bg-gray-200 dark:hover:bg-gray-700" role="alert">
            <span class="text-xs bg-primary-600 rounded-full text-white px-4 py-1.5 mr-3">New</span> 
            <span class="text-sm font-medium text-gray-800 dark:text-white">Lattice.go <strong>v1.0</strong> is released! See what's new</span>
            <svg class="ml-2 w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd"></path></svg>
        </a>
        <header class="mb-4 font-extrabold tracking-wide leading-none text-gray-900 md:text-5xl lg:text-6xl dark:text-white">
            Build Secure, Decentralized Web3 Apps on ZLattice
        </header>
        <p class="mb-8 text-lg font-medium tracking-tight font-sans text-gray-500 lg:text-xl sm:px-16 xl:px-48 dark:text-gray-400">
            A powerful SDK for ZLattice with multi-curve support. Deploy contracts, send transactions, and interact with smart contracts effortlessly.
        </p>
        <div class="flex flex-col mb-8 lg:mb-16 space-y-4 sm:flex-row sm:justify-center sm:space-y-0 sm:space-x-4">
            <div class="no-underline-important inline-flex justify-center items-center py-3 px-5 text-base font-medium text-center text-white rounded-lg bg-primary-700 hover:bg-primary-800 focus:ring-4 focus:ring-primary-300 dark:focus:ring-primary-900">
                <a class="flex items-center cursor-pointer !no-underline !text-white" href="/guide/getting-started/">
                    Get Started
                    <ArrowRightLong />
                </a>
            </div>
            <div class="no-underline-important inline-flex justify-center items-center py-3 px-5 text-base font-medium text-center text-gray-900 rounded-lg border border-gray-300 hover:bg-gray-100 focus:ring-4 focus:ring-gray-100 dark:text-white dark:border-gray-700 dark:hover:bg-gray-700 dark:focus:ring-gray-800">
                <a class="flex items-center !no-underline !text-gray-700 dark:!text-white cursor-pointer" href="/resources/ledger">
                    <Resources />
                    View Resources
                </a>
            </div>  
        </div>
    </div>
</section>
</template>

<script>
import ArrowRightLong from './icons/ArrowRightLong.vue';
import Resources from './icons/Resources.vue';

export default {
    name: 'Hero',
    components: {
        ArrowRightLong,
        Resources
    }
}
</script>

<style scoped>

</style>
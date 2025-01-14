<template>
  <div class="w-full py-12">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="text-center mb-6">
        <!-- Features badge -->
        <span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-indigo-100 text-indigo-600 dark:bg-indigo-900/50 dark:text-indigo-300">
          Features
        </span>
        
        <!-- Main title -->
        <span class="mt-6 text-4xl font-bold text-gray-900 dark:text-white block">
          Beautiful blockchain to grow smarter
        </span>
        
        <!-- Description text with better centering -->
        <div class="mt-6 flex justify-center">
          <p class="text-xl text-gray-600 dark:text-gray-400 max-w-3xl text-center">
            Powerful, self-serve blockchain and smart contract platform to help you build, deploy,
            <br class="hidden sm:block" />
            and scale your applications. Trusted by over 1,000 enterprises.
          </p>
        </div>
      </div>

      <div class="mt-20 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-x-8 gap-y-16">
        <div v-for="feature in features" :key="feature.title" class="relative">
          <div class="flex items-center justify-center mb-4">
            <div class="p-2 bg-indigo-50 dark:bg-indigo-900/50 rounded-lg">
              <svg class="w-6 h-6 text-indigo-600 dark:text-indigo-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="feature.icon"></path>
              </svg>
            </div>
          </div>
          <h3 class="text-xl font-semibold text-gray-900 dark:text-white text-center">{{ feature.title }}</h3>
          <p class="mt-3 text-base text-gray-500 dark:text-gray-400 text-center">
            {{ feature.description }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Features',
  data() {
    return {
      features: [
        {
          title: 'Self-Developed Blockchain',
          description: 'Self-developed consensus and encryption algorithms, perfectly compatible with secp256k1 and sm2p256v1, providing robust security guarantees.',
          icon: 'M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4'
        },
        {
          title: 'Asset Transfer',
          description: 'Fast and secure transfer functionality, supporting transactions of multiple asset types.',
          icon: 'M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4'
        },
        {
          title: 'EVM Compatible',
          description: 'Fully compatible with Ethereum Virtual Machine (EVM) smart contracts, seamlessly connecting with existing ecosystems.',
          icon: 'M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10'
        },
        {
          title: 'Built-in Contracts',
          description: 'Rich library of built-in contracts, including notary and other practical contracts for rapid business implementation.',
          icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z'
        },
        {
          title: 'High-Performance Architecture',
          description: 'Layered architecture design supporting thousands of transactions per second, meeting enterprise-level application needs.',
          icon: 'M13 10V3L4 14h7v7l9-11h-7z'
        },
        {
          title: 'Privacy Protection',
          description: 'Supports national cryptographic algorithms with multi-level data encryption and access control mechanisms.',
          icon: 'M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z'
        }
      ]
    }
  }
}
</script>

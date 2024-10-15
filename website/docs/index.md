---
# https://vitepress.dev/reference/default-theme-home-page
layout: home

hero:
  name: "Blockchain"
  text: "Lightweight Go library for integration with ZLattice clients."
  tagline: My great project tagline
  actions:
    - theme: brand
      text: Markdown Examples
      link: /markdown-examples
    - theme: alt
      text: API Examples
      link: /api-examples
---

<script setup>
  import Card from './components/Card.vue';
  import { VPTeamMembers } from 'vitepress/theme'

  const members = [
    {
      avatar: 'https://avatars.githubusercontent.com/u/114670506?v=4',
      name: 'Wenyang Lu',
      title: 'SDK',
      links: [
        { icon: 'github', link: 'https://github.com/wylu1037' },
        { icon: 'twitter', link: 'https://twitter.com' }
      ]
    },
    {
      avatar: 'https://www.github.com/yyx990803.png',
      name: 'Evan You',
      title: 'Creator',
      links: [
        { icon: 'github', link: 'https://github.com/yyx990803' },
        { icon: 'twitter', link: 'https://twitter.com/youyuxi' }
      ]
    },
    {
      avatar: 'https://www.github.com/yyx990803.png',
      name: 'Evan You',
      title: 'Creator',
      links: [
        { icon: 'github', link: 'https://github.com/yyx990803' },
        { icon: 'twitter', link: 'https://twitter.com/youyuxi' }
      ]
    },
    {
      avatar: 'https://www.github.com/yyx990803.png',
      name: 'Evan You',
      title: 'Creator',
      links: [
        { icon: 'github', link: 'https://github.com/yyx990803' },
        { icon: 'twitter', link: 'https://twitter.com/youyuxi' }
      ]
    },
  ]
</script>

<div class="flex justify-center text-5xl font-extrabold">
  <span class="bg-clip-text text-transparent bg-gradient-to-r from-pink-500 to-violet-500">
    Our Team
  </span>
</div>


<VPTeamMembers size="small" :members="members" />

<div class="flex justify-start my-10">
  <Card />
</div>

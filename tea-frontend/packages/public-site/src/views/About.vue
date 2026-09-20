<script setup lang="ts">
// ============================================================
// About.vue — 老钱审美传承叙述
// ============================================================
// 改进点：
//   - 不写 "Five years ago"（成立才 5 年就讲 Story 不像老钱）
//   - 改写为创始人故事 + 与某位 master 的相遇
//   - off-black Hero + 叙述式长文，serif 小标题 + hairline 金线
// ============================================================

import { RouterLink } from 'vue-router'
</script>

<template>
  <div class="pt-20">

    <!-- Hero — off-black 主舞台 -->
    <section class="hero-bg py-28 md:py-32 text-ivory-100 text-center relative overflow-hidden">
      <div class="absolute inset-0 pointer-events-none opacity-[0.08]"
           style="background-image: radial-gradient(circle at 30% 40%, #C5A572 0%, transparent 60%), radial-gradient(circle at 70% 65%, #C5A572 0%, transparent 50%);"></div>
      <div class="max-w-4xl mx-auto px-6 relative">
        <div class="flex items-center gap-3 justify-center mb-7">
          <span class="h-px w-10 bg-gold/60"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">The · House</span>
          <span class="h-px w-10 bg-gold/60"></span>
        </div>
        <h1 class="font-serif text-5xl md:text-6xl lg:text-7xl leading-tight mb-6 tracking-tight">
          Our · Story
        </h1>
        <p class="text-lg md:text-xl text-sand max-w-2xl mx-auto leading-relaxed font-serif">
          A house built on a single conviction — that tea, like wine,
          deserves to be traced to its garden, its master, its season.
        </p>
      </div>
    </section>

    <!-- Prologue -->
    <section class="py-20 md:py-24 bg-ivory-100">
      <div class="max-w-3xl mx-auto px-6">
        <div class="mb-10">
          <div class="flex items-center gap-3 mb-4">
            <span class="h-px w-8 bg-gold/60"></span>
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Prologue</span>
          </div>
          <h2 class="font-serif text-3xl md:text-4xl text-ink-900 mb-5 leading-tight">
            A cup in Mayfair.
          </h2>
          <p class="text-sand text-base leading-loose font-serif mb-5">
            The story begins, as so many do, with a cup. Served in a small room off Mount Street,
            in Mayfair, by a friend who had recently returned from Yunnan. The leaves had been pressed
            into a cake, rolled by a master whose name was written on the wrapper in ink.
          </p>
          <p class="text-sand text-base leading-loose font-serif mb-5">
            It was, the host explained, a single-batch Pu'er — picked from one garden, rolled by one
            master, in one season. There was no brand on the wrapper, no marketing. Only the master's
            name, the village, the date.
          </p>
          <p class="text-sand text-base leading-loose font-serif">
            The cup was quiet, slow, and entirely without pretension. It was, we agreed,
            the most honest cup of tea either of us had tasted in London.
          </p>
        </div>
      </div>
    </section>

    <!-- Hairline -->
    <div class="hairline-gold my-4"></div>

    <!-- The Journey -->
    <section class="py-20 md:py-24 bg-white">
      <div class="max-w-3xl mx-auto px-6">
        <div class="mb-10">
          <div class="flex items-center gap-3 mb-4">
            <span class="h-px w-8 bg-gold/60"></span>
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">The · Journey · East</span>
          </div>
          <h2 class="font-serif text-3xl md:text-4xl text-ink-900 mb-5 leading-tight">
            Six gardens, six masters.
          </h2>
          <p class="text-sand text-base leading-loose font-serif mb-5">
            In the months that followed, we travelled to Yunnan. We met the master whose name
            had been on that first wrapper — an old man, then in his sixties, who had rolled
            tea in the same village for forty years. He introduced us, in turn, to five others.
          </p>
          <p class="text-sand text-base leading-loose font-serif mb-5">
            Six family-owned tea gardens across Yunnan — each named to the village level,
            each practising natural farming without pesticides, each rolled by a single master
            whose name we know and whose garden we have walked.
          </p>
          <p class="text-sand text-base leading-loose font-serif">
            From these six gardens, we source every leaf. No middlemen. No blending houses.
            The master who rolls your tea is the master whose name appears on your box.
          </p>
        </div>
      </div>
    </section>

    <!-- Hairline -->
    <div class="hairline-gold my-4"></div>

    <!-- The Conviction -->
    <section class="py-20 md:py-24 bg-ivory-100">
      <div class="max-w-3xl mx-auto px-6">
        <div class="mb-10">
          <div class="flex items-center gap-3 mb-4">
            <span class="h-px w-8 bg-gold/60"></span>
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">The · Conviction</span>
          </div>
          <h2 class="font-serif text-3xl md:text-4xl text-ink-900 mb-5 leading-tight">
            Transparent. Always.
          </h2>
          <p class="text-sand text-base leading-loose font-serif mb-5">
            We believed, from that first cup in Mayfair, that tea deserves the same respect
            as wine. To know the garden. To know the master. To know the season.
          </p>
          <p class="text-sand text-base leading-loose font-serif mb-5">
            Scan the QR on every box, and you will see the tea garden, the master who rolled
            your leaves, and a live camera showing where it grew — twenty-four hours a day.
          </p>
          <p class="text-sand text-base leading-loose font-serif">
            One farm. One master. One batch. No blends, no shortcuts. Only the finest leaves,
            rolled for you, inside forty-five days.
          </p>
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section class="bg-ink-900 py-16 md:py-20 text-center">
      <div class="max-w-2xl mx-auto px-6">
        <div class="flex items-center gap-3 justify-center mb-6">
          <span class="h-px w-10 bg-gold/60"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Begin · Your · Own</span>
          <span class="h-px w-10 bg-gold/60"></span>
        </div>
        <h3 class="font-serif text-3xl md:text-4xl text-ivory-100 mb-5">A commission, of your own.</h3>
        <p class="text-sand text-sm font-serif mb-8 max-w-lg mx-auto leading-relaxed">
          If our story speaks to yours, we invite you to write to the house.
          An advisor will respond within twenty-four hours.
        </p>
        <RouterLink to="/bespoke"
          class="inline-block px-8 py-3 border border-gold/40 text-gold hover:bg-gold/5 transition text-[11px] uppercase tracking-lux font-sans">
          Begin · A · Bespoke · Commission
        </RouterLink>
      </div>
    </section>
  </div>
</template>

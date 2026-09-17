<script setup lang="ts">
const mountains = [
  { name: '冰岛老寨', region: '临沧 · 勐库', altitude: '2000m+', specialty: '冰糖甜感', season: 'Apr—May', master: '李师傅', story: '1999年古树茶价格复苏，冰岛老寨成为普洱界的"圣杯"。云雾缭绕，茶树沿坡而种。' },
  { name: '班章村', region: '西双版纳 · 勐海', altitude: '1700m+', specialty: '霸气回甘', season: 'Mar—Apr', master: '张师傅', story: '"班章为王"——厚重的山野气息，浓郁的喉韵，是最具识别度的普洱之一。' },
  { name: '景迈山', region: '普洱 · 澜沧', altitude: '1400m+', specialty: '兰花蜜香', season: 'Apr—May', master: '王师傅', story: '千年万亩古茶园，2010 年入选全球重要农业文化遗产。' },
  { name: '老班章', region: '勐海 · 布朗山', altitude: '1800m', specialty: '苦尽甘来', season: 'Apr', master: '赵师傅', story: '原始森林中的古树群落，强烈的山野气，头春茶千金难求。' },
  { name: '南糯山', region: '勐海 · 格朗和', altitude: '1600m+', specialty: '高扬花香', season: 'Mar—Apr', master: '陈师傅', story: '傣族、哈尼族世代种茶之地，乔木茶林连绵不绝。' },
  { name: '忙肺村', region: '临沧 · 永德', altitude: '1800m', specialty: '清凉兰香', season: 'Apr—May', master: '周师傅', story: '大雪山深处，小种茶树天然有机，茶气清新独特。' },
]
</script>

<template>
  <div class="pt-20">
    <section class="bg-tea-900 text-tea-50 py-24">
      <div class="max-w-6xl mx-auto px-6">
        <span class="text-xs tracking-widest uppercase text-tea-400">茶山上 — Where It All Starts</span>
        <h1 class="font-serif text-5xl md:text-6xl mt-4 mb-6">The Mountains<br/>We Source From.</h1>
        <p class="text-lg text-tea-200 max-w-2xl leading-relaxed">Yunnan's cloud mountains produce 98% of the world's Pu'er. Every farm we work with is GPS-tagged and video-monitored 24/7.</p>
      </div>
    </section>

    <section class="py-24 bg-tea-50">
      <div class="max-w-7xl mx-auto px-6">
        <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
          <article v-for="m in mountains" :key="m.name" class="bg-white rounded-2xl overflow-hidden card-hover border border-tea-100">
            <div class="h-40 bg-gradient-to-br from-tea-700 to-tea-400 flex items-center justify-center text-6xl">⛰️</div>
            <div class="p-6">
              <h3 class="font-serif text-xl text-tea-900">{{ m.name }}</h3>
              <p class="text-xs text-tea-500 mt-1 mb-4">{{ m.region }} · {{ m.altitude }}</p>
              <p class="text-sm text-tea-700 leading-relaxed mb-4">{{ m.story }}</p>
              <div class="flex flex-wrap gap-2 text-xs">
                <span class="px-2 py-1 bg-tea-100 text-tea-800 rounded">✨ {{ m.specialty }}</span>
                <span class="px-2 py-1 bg-tea-100 text-tea-800 rounded">🌱 {{ m.season }}</span>
                <span class="px-2 py-1 bg-tea-100 text-tea-800 rounded">👤 Master {{ m.master }}</span>
              </div>
            </div>
          </article>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const gardens = ref([
  {
    name: '云南省临沧市临翔区邦东乡曼岗村茶园',
    region: '临沧 · 云南',
    master: '王师傅',
    tea: '古树普洱生茶',
    harvest: '4 月中旬（清明后）',
    elevation: '约 1800m',
    desc: '邦东乡地处临沧东北，海拔高、昼夜温差大，是古树普洱的核心产区。曼岗村的茶树龄多在 200-300 年之间，茶气足、回甘快。我们与该村合作社长期合作，所有茶叶手工采摘。',
    stream_url: 'rtmp://localhost:1935/slow/mangang',
    status: 'live',
  },
  {
    name: '广东省潮州市潮安区凤凰镇大乌岽村茶园',
    region: '潮州 · 广东',
    master: '李师傅',
    tea: '凤凰单丛',
    harvest: '3-4 月（春茶）',
    elevation: '约 1300m',
    desc: '凤凰单丛以香型丰富著称，蜜兰香、鸭屎香、黄枝香是最受欢迎的三个品系。大乌岽村海拔较高，终年云雾缭绕，茶香高扬持久。',
    stream_url: 'rtmp://localhost:1935/slow/dawudong',
    status: 'live',
  },
  {
    name: '云南省普洱市澜沧拉祜族自治县惠民镇景迈村茶园',
    region: '澜沧 · 云南',
    master: '岩师父（拉祜族）',
    tea: '普洱生茶',
    harvest: '4-5 月',
    elevation: '约 1400m',
    desc: '景迈村是拉祜族和傣族世代居住的古老茶村，种茶历史超过 1300 年。茶园与原始森林伴生，生态环境极佳，茶叶带有独特的兰花香。',
    stream_url: 'rtmp://localhost:1935/slow/jingmai',
    status: 'idle',
  },
  {
    name: '云南省西双版纳州勐海县布朗山乡班章村茶园',
    region: '勐海 · 云南',
    master: '张师傅',
    tea: '古树普洱生茶',
    harvest: '4 月',
    elevation: '约 1600m',
    desc: '班章村位于布朗山深处，茶树根系深扎岩缝，茶气刚猛、厚重饱满，是普洱茶爱好者心中的标杆产区。我们只与村合作社合作，保护古树资源。',
    stream_url: 'rtmp://localhost:1935/slow/banzhang',
    status: 'live',
  },
  {
    name: '云南省普洱市澜沧拉祜族自治县惠民镇翁基村茶园',
    region: '澜沧 · 云南',
    master: '白师父（布朗族）',
    tea: '普洱生茶',
    harvest: '4-5 月',
    elevation: '约 1500m',
    desc: '翁基村是布朗族的千年古寨，家家户户种茶制茶，保留着传统的手工杀青工艺。茶叶带有独特的甜香和木质调。',
    stream_url: 'rtmp://localhost:1935/slow/wengji',
    status: 'idle',
  },
  {
    name: '云南省临沧市临翔区南美乡多依村茶园',
    region: '临沧 · 云南',
    master: '李师傅',
    tea: '古树普洱生茶',
    harvest: '4 月下旬',
    elevation: '约 1700m',
    desc: '多依村地处南美乡最高处，终年云雾笼罩，茶树生长缓慢，叶质肥厚。茶叶带有明显的山野气和花蜜香。',
    stream_url: 'rtmp://localhost:1935/slow/duoyi',
    status: 'idle',
  },
])
</script>

<template>
  <div class="pt-20 min-h-screen">
    <section class="bg-tea-900 text-tea-50 py-20">
      <div class="max-w-6xl mx-auto px-6">
        <span class="text-xs tracking-[0.3em] text-tea-300 uppercase">Tea Garden Stories</span>
        <h1 class="font-serif text-5xl md:text-6xl mb-6">Six Tea Gardens.<br/>One Promise.</h1>
        <p class="text-lg text-tea-300 max-w-2xl">Every batch is traced back to the village where it was grown. No GPS coordinates — out of respect for local communities. Just the village name, the master who picked it, and the story of the land.</p>
      </div>
    </section>

    <section class="py-16 bg-tea-50">
      <div class="max-w-7xl mx-auto px-6">
        <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div v-for="g in gardens" :key="g.name" class="bg-white rounded-2xl overflow-hidden border border-tea-100 card-hover">
            <div class="relative aspect-video bg-tea-900 flex items-center justify-center text-tea-200">
              <div class="text-center">
                <div class="text-5xl mb-2">🍃</div>
                <div class="text-xs font-mono text-tea-400">{{ g.stream_url }}</div>
              </div>
              <div v-if="g.status === 'live'" class="absolute top-3 left-3 flex items-center gap-2 px-3 py-1 bg-red-600 text-white text-xs font-medium rounded-full">
                <span class="w-1.5 h-1.5 rounded-full bg-white animate-pulse"></span> LIVE
              </div>
              <div v-else class="absolute top-3 left-3 flex items-center gap-2 px-3 py-1 bg-tea-600 text-white text-xs font-medium rounded-full">
                Offline
              </div>
            </div>
            <div class="p-5 space-y-3">
              <div>
                <h3 class="font-serif text-lg text-tea-900 leading-snug">{{ g.name }}</h3>
                <p class="text-xs text-tea-500 mt-1">{{ g.region }} · 海拔 {{ g.elevation }}</p>
              </div>
              <p class="text-sm text-tea-600 leading-relaxed">{{ g.desc }}</p>
              <div class="flex gap-4 text-xs pt-3 border-t border-tea-100">
                <div><div class="text-tea-400">茶师</div><div class="text-tea-900 font-medium">{{ g.master }}</div></div>
                <div><div class="text-tea-400">品类</div><div class="text-tea-900 font-medium">{{ g.tea }}</div></div>
                <div><div class="text-tea-400">采摘季</div><div class="text-tea-900 font-medium">{{ g.harvest }}</div></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

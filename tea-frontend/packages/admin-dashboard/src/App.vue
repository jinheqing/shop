<script setup lang="ts">
import { ElContainer, ElHeader, ElAside, ElMain, ElMenu } from 'element-plus'
import { Menu } from '@element-plus/icons-vue'
</script>

<template>
  <ElContainer class="h-screen">
    <ElHeader class="flex items-center justify-between border-b border-gray-200 bg-white">
      <div class="flex items-center gap-3">
        <el-icon class="text-mattegold text-xl"><Menu /></el-icon>
        <span class="font-semibold text-lg text-ink">Tea System · 管理后台</span>
      </div>
      <div class="text-sm text-gray-400">v0.1.0 · Step 1 scaffolding</div>
    </ElHeader>
    <ElContainer>
      <ElAside width="240px" class="border-r border-gray-200 bg-gray-50">
        <ElMenu :default-active="'dashboard'" class="border-r-0 bg-gray-50">
          <ElMenuItem index="dashboard">Dashboard</ElMenuItem>
          <ElMenuItem index="im">IM 会话</ElMenuItem>
          <ElMenuItem index="orders">订单管理</ElMenuItem>
          <ElMenuItem index="quotes">定制报价</ElMenuItem>
          <ElMenuItem index="live">🎥 直播管理</ElMenuItem>
          <ElMenuItem index="live-slow">🎥 慢直播管理</ElMenuItem>
          <ElMenuItem index="nodes">节点管理</ElMenuItem>
          <ElMenuItem index="settings">设置</ElMenuItem>
        </ElMenu>
      </ElAside>
      <ElMain class="p-8 bg-white">
        <div class="text-gray-400">Admin dashboard scaffolding — Step 1 complete.</div>
      </ElMain>
    </ElContainer>
  </ElContainer>
</template>

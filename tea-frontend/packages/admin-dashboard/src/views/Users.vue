<template>
  <div>
    <el-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-bold text-lg">{{ t('users.title') }}</span>
          <el-button :icon="Refresh" @click="load">{{ t('users.reload') }}</el-button>
        </div>
      </template>
      <el-table :data="items" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" :label="t('users.name')" width="140" />
        <el-table-column prop="email" :label="t('users.email')" width="220" />
        <el-table-column :label="t('users.location')" width="200">
          <template #default="{ row }">
            <div v-if="row.last_login_city || row.last_login_country" class="flex items-center gap-1">
              <span>{{ getFlag(row.last_login_country_code) }}</span>
              <span>{{ row.last_login_city }}{{ row.last_login_region ? ', ' + row.last_login_region : '' }}</span>
              <span v-if="row.last_login_country" class="text-gray-400 text-xs">({{ row.last_login_country }})</span>
            </div>
            <span v-else class="text-gray-400">—</span>
          </template>
        </el-table-column>
        <el-table-column prop="last_login_ip" :label="t('users.ip')" width="140">
          <template #default="{ row }">
            <span v-if="row.last_login_ip" class="font-mono text-xs">{{ row.last_login_ip }}</span>
            <span v-else class="text-gray-400">—</span>
          </template>
        </el-table-column>
        <el-table-column prop="is_email_verified" :label="t('users.verified')" width="90">
          <template #default="{ row }">
            <el-tag v-if="row.is_email_verified" type="success">{{ t('users.yes') }}</el-tag>
            <el-tag v-else type="warning">{{ t('users.no') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_login_at" :label="t('users.last_login')" width="170">
          <template #default="{ row }">
            <span v-if="row.last_login_at" class="text-sm">{{ formatTime(row.last_login_at) }}</span>
            <span v-else class="text-gray-400">—</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="t('users.signed_up')" width="170" />
        <el-table-column :label="t('users.actions')" width="140">
          <template #default="{ row }">
            <el-button size="small" @click="viewOrders(row)">{{ t('users.orders') }}</el-button>
            <el-button size="small" type="danger" @click="exportData(row)">{{ t('users.export') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { api } from '../api/client'

const { t, locale } = useI18n()
const items = ref<any[]>([])

async function load() {
  const res: any = await api.get('/users')
  items.value = res?.items || res || []
}

function viewOrders(row: any) {
  window.open(`/#/orders?user_id=${row.id}`, '_blank')
}

async function exportData(row: any) {
  ElMessage.info(`Would trigger DSAR export for user ${row.id}`)
}

// 将国家代码转为旗帜 emoji
function getFlag(code?: string): string {
  if (!code || code.length !== 2) return ''
  const A = 0x1F1E6
  const a = 'A'.charCodeAt(0)
  return String.fromCodePoint(A + (code.charCodeAt(0) - a), A + (code.charCodeAt(1) - a))
}

function formatTime(ts: string): string {
  if (!ts) return '—'
  try {
    return new Date(ts).toLocaleString(locale.value === 'zh' ? 'zh-CN' : 'en-GB', {
      year: 'numeric', month: '2-digit', day: '2-digit',
      hour: '2-digit', minute: '2-digit'
    })
  } catch {
    return ts
  }
}

onMounted(load)
</script>

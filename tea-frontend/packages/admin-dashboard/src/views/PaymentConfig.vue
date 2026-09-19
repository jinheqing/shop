<template>
  <div>
    <el-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-bold text-lg">System — Payment Configuration</span>
          <div>
            <el-tag v-if="saving" type="warning">Saving…</el-tag>
            <el-tag v-else-if="savedAt" type="success">Saved {{ savedAt }}</el-tag>
          </div>
        </div>
      </template>

      <el-alert v-if="!loaded" title="Loading from /system/config/payment_config…" type="info" show-icon :closable="false" class="mb-4" />

      <el-form v-else :model="cfg" label-width="200px" class="max-w-3xl">
        <!-- 2Checkout -->
        <el-divider content-position="left">2Checkout (Primary Gateway)</el-divider>
        <el-form-item label="Account #">
          <el-input v-model="cfg.twocheckout_account_id" placeholder="1234567890" />
          <div class="text-xs text-gray-400 mt-1">NOT the seller ID — use the numeric account identifier</div>
        </el-form-item>
        <el-form-item label="Secret Key">
          <el-input v-model="cfg.twocheckout_secret_key" type="password" show-password placeholder="Encrypted at rest" />
        </el-form-item>
        <el-form-item label="3DS / SCA">
          <el-switch v-model="cfg.twocheckout_3ds_enabled" />
          <span class="ml-2 text-xs text-gray-500">Strong Customer Authentication (EU required)</span>
        </el-form-item>
        <el-form-item label="Mode">
          <el-radio-group v-model="cfg.twocheckout_mode">
            <el-radio value="sandbox">Sandbox</el-radio>
            <el-radio value="live">Live</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="Webhook URL">
          <el-input v-model="cfg.twocheckout_webhook_url" placeholder="/api/v1/webhooks/2checkout" />
        </el-form-item>

        <!-- PayPal -->
        <el-divider content-position="left">PayPal (Optional Alternative)</el-divider>
        <el-form-item label="Client ID">
          <el-input v-model="cfg.paypal_client_id" />
        </el-form-item>
        <el-form-item label="Client Secret">
          <el-input v-model="cfg.paypal_client_secret" type="password" show-password />
        </el-form-item>
        <el-form-item label="Mode">
          <el-radio-group v-model="cfg.paypal_mode">
            <el-radio value="sandbox">Sandbox</el-radio>
            <el-radio value="live">Live</el-radio>
          </el-radio-group>
        </el-form-item>

        <!-- 通用 -->
        <el-divider content-position="left">General</el-divider>
        <el-form-item label="Currency">
          <el-select v-model="cfg.currency">
            <el-option label="GBP — Pound Sterling" value="GBP" />
            <el-option label="EUR — Euro" value="EUR" />
            <el-option label="USD — US Dollar" value="USD" />
          </el-select>
        </el-form-item>
        <el-form-item label="VAT Rate (%)">
          <el-input-number v-model="cfg.vat_rate" :precision="2" :step="0.5" :min="0" :max="50" />
        </el-form-item>

        <!-- 推荐人自动规则 -->
        <el-divider content-position="left">Referral Automation</el-divider>
        <el-alert type="info" :closable="false" show-icon class="mb-3">
          推荐人系统规则 — 无需重新上线，改完立即生效。默认：£100 / 累计 3 个
        </el-alert>
        <el-form-item label="Friend Order Threshold (£)">
          <el-input-number v-model="cfg.referral_min_order_amount" :precision="0" :step="10" :min="0" />
          <div class="text-xs text-gray-400 mt-1">朋友订单金额 ≥ 这个数才触发推荐奖励（礼物 + 进组检查）。设 0 = 所有订单都触发</div>
        </el-form-item>
        <el-form-item label="Auto-Invite to KOL Group">
          <el-input-number v-model="cfg.referral_count_to_vip" :precision="0" :step="1" :min="1" :max="50" />
          <div class="text-xs text-gray-400 mt-1">推荐人累计有多少个"有效推荐"后，系统自动把他拉进 KOL 伙伴用户组（闭门品鉴会、限量茶先购、茶园参观）</div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="saving" @click="save">Save All</el-button>
          <el-button @click="load">Reload</el-button>
          <el-button @click="testWebhook" :disabled="!cfg.twocheckout_webhook_url">Test 2Checkout Webhook</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '../api/client'

const loaded = ref(false)
const saving = ref(false)
const savedAt = ref('')
const cfg = ref<any>({
  twocheckout_account_id: '',
  twocheckout_secret_key: '',
  twocheckout_3ds_enabled: true,
  twocheckout_mode: 'sandbox',
  twocheckout_webhook_url: '/api/v1/webhooks/2checkout',
  paypal_client_id: '',
  paypal_client_secret: '',
  paypal_mode: 'sandbox',
  currency: 'GBP',
  vat_rate: 20,
  referral_min_order_amount: 100,
  referral_count_to_vip: 3
})

async function load() {
  try {
    const res: any = await api.get('/system/config/payment_config')
    if (res?.value && typeof res.value === 'object') {
      cfg.value = { ...cfg.value, ...res.value }
    }
    loaded.value = true
  } catch (e: any) {
    ElMessage.error('Load failed: ' + (e?.message || e))
  }
}

async function save() {
  saving.value = true
  try {
    await api.put('/system/config/payment_config', { value: cfg.value })
    ElMessage.success('Saved successfully')
    savedAt.value = new Date().toLocaleTimeString()
  } catch (e: any) {
    ElMessage.error('Save failed: ' + (e?.message || e))
  } finally {
    saving.value = false
  }
}

function testWebhook() {
  ElMessage.info('Webhook test endpoint would POST a mock event — not wired yet')
}

onMounted(load)
</script>

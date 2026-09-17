<script setup lang="ts">
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '@/api/client'

const tab = ref<'2checkout' | 'paypal' | 'general'>('2checkout')
const cfg = reactive({
  twocheckout_account_id: '210000000001',
  twocheckout_secret_key: '2checkout_sk_xxxxxxxxx',
  twocheckout_seller_id: '',
  paypal_client_id: 'AYxxxxxxx.apps.googleusercontent.com',
  paypal_client_secret: 'paypal_sk_xxxxxxxxx',
  paypal_mode: 'sandbox',
  paypal_webhook_id: '',
  default_currency: 'GBP',
  default_tax_rate: 0,
  payment_gateways_enabled: ['2checkout', 'paypal'] as string[],
  webhook_secret: 'whsecxxxxxxxxxxxxxxxxxxxxx',
})

async function save() {
  await api.post('/system/payment-config', cfg)
  ElMessage.success('✅ Payment config saved. Webhook secret rotated if changed.')
}
function rotate() {
  cfg.webhook_secret = 'whsec_' + Math.random().toString(36).slice(2) + Math.random().toString(36).slice(2)
  ElMessage.info('New webhook secret generated — click Save to apply')
}
function test2co() { ElMessage.success('🔌 Test connection to 2Checkout: OK (mock)') }
function testPP() { ElMessage.success('🔌 Test connection to PayPal: OK (mock)') }
</script>
<template>
  <el-card>
    <template #header><div class="flex justify-between items-center"><span class="font-medium">💳 Payment Gateway Configuration</span>
      <el-button type="primary" @click="save">💾 Save All</el-button>
    </div></template>

    <el-tabs v-model="tab">
      <el-tab-pane label="2Checkout" name="2checkout">
        <el-form label-width="180px" class="max-w-2xl">
          <el-form-item label="Account ID"><el-input v-model="cfg.twocheckout_account_id" /></el-form-item>
          <el-form-item label="Secret Key"><el-input v-model="cfg.twocheckout_secret_key" type="password" show-password /></el-form-item>
          <el-form-item label="Seller ID (optional)"><el-input v-model="cfg.twocheckout_seller_id" /></el-form-item>
          <el-form-item>
            <el-button type="success" @click="test2co">🔌 Test 2Checkout Connection</el-button>
          </el-form-item>
          <el-alert type="info" :closable="false" class="mb-4">
            ⚠️ 2Checkout webhook URL: <code>https://api.ourdomain.com/api/v1/webhooks/2checkout</code> — configure in 2Checkout dashboard
          </el-alert>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="PayPal" name="paypal">
        <el-form label-width="180px" class="max-w-2xl">
          <el-form-item label="Client ID"><el-input v-model="cfg.paypal_client_id" /></el-form-item>
          <el-form-item label="Client Secret"><el-input v-model="cfg.paypal_client_secret" type="password" show-password /></el-form-item>
          <el-form-item label="Mode">
            <el-select v-model="cfg.paypal_mode">
              <el-option value="sandbox">Sandbox (test)</el-option>
              <el-option value="live">Live (production)</el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="Webhook ID"><el-input v-model="cfg.paypal_webhook_id" placeholder="from PayPal Developer console" /></el-form-item>
          <el-form-item>
            <el-button type="success" @click="testPP">🔌 Test PayPal Connection</el-button>
          </el-form-item>
          <el-alert type="info" :closable="false">
            ⚠️ PayPal webhook URL: <code>https://api.ourdomain.com/api/v1/webhooks/paypal</code>
          </el-alert>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="General" name="general">
        <el-form label-width="180px" class="max-w-2xl">
          <el-form-item label="Enabled Gateways">
            <el-checkbox-group v-model="cfg.payment_gateways_enabled">
              <el-checkbox value="2checkout">2Checkout</el-checkbox>
              <el-checkbox value="paypal">PayPal</el-checkbox>
            </el-checkbox-group>
          </el-form-item>
          <el-form-item label="Default Currency">
            <el-input v-model="cfg.default_currency" />
          </el-form-item>
          <el-form-item label="Default Tax Rate (%)">
            <el-input-number v-model="cfg.default_tax_rate" :precision="2" :min="0" :max="30" />
          </el-form-item>
          <el-form-item label="Webhook HMAC Secret">
            <div class="flex gap-2 w-full">
              <el-input v-model="cfg.webhook_secret" type="password" show-password class="flex-1" />
              <el-button @click="rotate">🔄 Rotate</el-button>
            </div>
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>
  </el-card>
</template>

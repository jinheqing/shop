<template>
  <div class="space-y-4">
    <el-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-bold text-lg">System — Runtime Configuration</span>
          <el-tag type="info" effect="plain">Changes require restart to take effect</el-tag>
        </div>
      </template>

      <el-alert
        title="All values below are persisted to the site_contents table (page_key = 'system_config')."
        type="info"
        :closable="false"
        class="mb-4"
      />

      <!-- 🔗 Quick links to sibling config pages -->
      <div class="flex gap-2 mb-4">
        <el-link href="#/payment-config" :underline="false">→ Payment Config (2Checkout / PayPal)</el-link>
        <el-link href="#/site-content" :underline="false">→ CMS (Public Site Pages)</el-link>
      </div>
    </el-card>

    <!-- Group 1: App / Branding -->
    <el-card>
      <template #header><b>① App &amp; Branding</b><span class="ml-2 text-xs text-gray-400">section = <code>app</code></span></template>
      <el-form :model="app" label-width="180px" class="max-w-3xl">
        <el-form-item label="Public Domain">
          <el-input v-model="app.domain" placeholder="https://ukteahouse.co.uk" />
          <div class="text-xs text-gray-400 mt-1">Used for magic-link URLs, OG meta, email From</div>
        </el-form-item>
        <el-form-item label="Service Name">
          <el-input v-model="app.service_name" placeholder="UK Tea House" />
        </el-form-item>
        <el-form-item label="Contact Email">
          <el-input v-model="app.contact_email" placeholder="tea@ukteahouse.co.uk" />
          <div class="text-xs text-gray-400 mt-1">Shown on GDPR / privacy pages</div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="saving.app" @click="save('app', app)">Save App</el-button>
          <el-tag v-if="savedAt.app" type="success" class="ml-2">Saved {{ savedAt.app }}</el-tag>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Group 2: Mail Service -->
    <el-card>
      <template #header><b>② Mail Service</b><span class="ml-2 text-xs text-gray-400">section = <code>mail</code></span></template>
      <el-form :model="mail" label-width="180px" class="max-w-3xl">
        <el-form-item label="SMTP Host">
          <el-input v-model="mail.smtp_host" placeholder="smtp.mailgun.org" />
        </el-form-item>
        <el-form-item label="From Address">
          <el-input v-model="mail.from_addr" placeholder="tea@ukteahouse.co.uk" />
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="mail.api_key" type="password" show-password placeholder="Mailgun / Resend API key" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="saving.mail" @click="save('mail', mail)">Save Mail</el-button>
          <el-tag v-if="savedAt.mail" type="success" class="ml-2">Saved {{ savedAt.mail }}</el-tag>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Group 3: Translate Engine -->
    <el-card>
      <template #header><b>③ Translate Engine</b><span class="ml-2 text-xs text-gray-400">section = <code>translate</code></span></template>
      <el-form :model="translate" label-width="200px" class="max-w-3xl">
        <el-form-item label="Engine Base URL">
          <el-input v-model="translate.service_url" placeholder="http://translate-engine:8090" />
          <div class="text-xs text-gray-400 mt-1">Default: http://localhost:8090 — must expose /translate/text endpoint</div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="saving.translate" @click="save('translate', translate)">Save Translate</el-button>
          <el-tag v-if="savedAt.translate" type="success" class="ml-2">Saved {{ savedAt.translate }}</el-tag>
          <el-button class="ml-2" @click="testTranslate">Test Connection</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Group 4: Media Nodes / WireGuard -->
    <el-card>
      <template #header><b>④ Media Nodes / WireGuard</b><span class="ml-2 text-xs text-gray-400">section = <code>nodes</code></span></template>
      <el-form :model="nodes" label-width="200px" class="max-w-3xl">
        <el-form-item label="WG IP Start">
          <el-input v-model="nodes.wg_ip_start" placeholder="10.10.0.10" />
          <div class="text-xs text-gray-400 mt-1">Used when creating new Media Nodes</div>
        </el-form-item>
        <el-form-item label="WG Prefix Length">
          <el-input v-model="nodes.wg_ip_prefix" placeholder="24" />
          <div class="text-xs text-gray-400 mt-1">CIDR bits: 24 = 255.255.255.0</div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="saving.nodes" @click="save('nodes', nodes)">Save Nodes</el-button>
          <el-tag v-if="savedAt.nodes" type="success" class="ml-2">Saved {{ savedAt.nodes }}</el-tag>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Group 5: Server Metadata -->
    <el-card>
      <template #header><b>⑤ Server Metadata</b><span class="ml-2 text-xs text-gray-400">section = <code>server</code></span></template>
      <el-form :model="server" label-width="180px" class="max-w-3xl">
        <el-form-item label="API Version">
          <el-input v-model="server.version" placeholder="0.5.0" />
          <div class="text-xs text-gray-400 mt-1">Returned by GET /health and GET /</div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="saving.server" @click="save('server', server)">Save Server</el-button>
          <el-tag v-if="savedAt.server" type="success" class="ml-2">Saved {{ savedAt.server }}</el-tag>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Group 6: Runtime Status (read-only) -->
    <el-card>
      <template #header><b>⑥ Runtime Status</b> <span class="text-xs text-gray-400">(read-only)</span></template>
      <div class="grid grid-cols-2 gap-4">
        <div>Config entries in DB: <b>{{ knownKeys.length }}</b></div>
        <div>Current loaded version: <b>{{ server.version || '(not set — using env default)' }}</b></div>
        <div>Service name: <b>{{ app.service_name }}</b></div>
        <div>Magic-link domain: <b>{{ app.domain || '(not set)' }}</b></div>
      </div>
      <el-button class="mt-3" @click="loadAll">Reload All from DB</el-button>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '../api/client'

// ——— 预置 schema：每个 key 对应一个默认 JSON body ———
type Section = Record<string, any>
const knownKeys = ref<string[]>(['app', 'mail', 'translate', 'nodes', 'server'])

const app      = reactive<Section>({ domain: '', service_name: '', contact_email: '' })
const mail     = reactive<Section>({ smtp_host: '', from_addr: '', api_key: '' })
const translate = reactive<Section>({ service_url: '' })
const nodes    = reactive<Section>({ wg_ip_start: '', wg_ip_prefix: '' })
const server   = reactive<Section>({ version: '' })

const saving = reactive<Record<string, boolean>>({})
const savedAt = reactive<Record<string, string>>({})

const sections: Record<string, Section> = { app, mail, translate, nodes, server }

async function loadAll() {
  for (const key of knownKeys.value) {
    try {
      const res: any = await api.get(`/system/config/${key}`)
      if (res?.value && typeof res.value === 'object') {
        Object.assign(sections[key], res.value)
      }
    } catch { /* 空 ok */ }
  }
  // 顺带拉一下 List，显示实际已存 key 数量
  try {
    const list: any = await api.get('/system/config')
    if (Array.isArray(list?.items)) {
      knownKeys.value = list.items.map((i: any) => i.key).length ? list.items.map((i: any) => i.key) : knownKeys.value
    }
  } catch { /* ignore */ }
  ElMessage.success('All sections reloaded')
}

async function save(key: string, body: Section) {
  saving[key] = true
  try {
    await api.put(`/system/config/${key}`, { value: body })
    ElMessage.success(`Saved '${key}' → DB (restart required)`)
    savedAt[key] = new Date().toLocaleTimeString()
  } catch (e: any) {
    ElMessage.error(`Save failed: ${e?.message || e}`)
  } finally {
    saving[key] = false
  }
}

function testTranslate() {
  ElMessage.info('Test endpoint not wired — enable TranslateEngine to verify')
}

onMounted(loadAll)
</script>

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
        <el-form-item label="Provider">
          <el-select v-model="mail.provider" placeholder="Select provider">
            <el-option label="AokSend (国内)" value="aoksend" />
            <el-option label="Mailgun (海外)" value="mailgun" />
          </el-select>
          <div class="text-xs text-gray-400 mt-1">选择邮件服务商，AokSend 适合国内用户，Mailgun 适合海外</div>
        </el-form-item>
        <el-form-item label="From Address">
          <el-input v-model="mail.from_addr" placeholder="noreply@tea.7758521.sbs" />
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="mail.api_key" type="password" show-password placeholder="AokSend / Mailgun API key" />
          <div class="text-xs text-gray-400 mt-1">在对应服务商后台获取 API Key</div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="saving.mail" @click="save('mail', mail)">Save Mail</el-button>
          <el-tag v-if="savedAt.mail" type="success" class="ml-2">Saved {{ savedAt.mail }}</el-tag>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Group 2.5: OBS Push/Pull Stream -->
    <el-card>
      <template #header><b>③ OBS Push/Pull Stream</b><span class="ml-2 text-xs text-gray-400">section = <code>obs</code></span></template>
      <el-alert title="配置 OBS 直播推流和拉流（播放）域名。留空则使用默认 MediaMTX 本地地址。" type="info" :closable="false" class="mb-3" />
      <el-form :model="obs" label-width="180px" class="max-w-3xl">
        <el-form-item label="Push Domain">
          <el-input v-model="obs.push_domain" placeholder="rtmp://38.76.188.92:1935/live" />
          <div class="text-xs text-gray-400 mt-1">OBS Studio 推流地址，格式: rtmp://IP:PORT/live</div>
        </el-form-item>
        <el-form-item label="Pull Domain">
          <el-input v-model="obs.pull_domain" placeholder="https://tea.7758521.sbs/webrtc/live" />
          <div class="text-xs text-gray-400 mt-1">WebRTC 播放地址，格式: https://域名/webrtc/live</div>
        </el-form-item>
        <el-form-item label="RTMP Port">
          <el-input-number v-model="obs.rtmp_port" :min="1" :max="65535" placeholder="1935" />
          <div class="text-xs text-gray-400 mt-1">默认 1935，一般不需要修改</div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="saving.obs" @click="save('obs', obs)">Save OBS</el-button>
          <el-tag v-if="savedAt.obs" type="success" class="ml-2">Saved {{ savedAt.obs }}</el-tag>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Group 3: Translate Engine -->
    <el-card>
      <template #header><b>④ Translate Engine</b><span class="ml-2 text-xs text-gray-400">section = <code>translate</code></span></template>
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
      <template #header><b>⑤ Media Nodes / WireGuard</b><span class="ml-2 text-xs text-gray-400">section = <code>nodes</code></span></template>
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
      <template #header><b>⑥ Server Metadata</b><span class="ml-2 text-xs text-gray-400">section = <code>server</code></span></template>
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

    <!-- Group 6: LiveKit SFU -->
    <el-card>
      <template #header><b>⑦ LiveKit SFU</b><span class="ml-2 text-xs text-gray-400">section = <code>livekit</code></span></template>
      <el-alert title="⚠️ LIVEKIT_API_KEY / LIVEKIT_API_SECRET are production secrets. Be careful editing here — values are stored in plaintext in site_contents (JSONB)." type="warning" show-icon :closable="false" class="mb-3" />
      <el-form :model="livekit" label-width="200px" class="max-w-3xl">
        <el-form-item label="SFU WebSocket URL">
          <el-input v-model="livekit.url" placeholder="http://livekit:7880" />
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="livekit.api_key" placeholder="livekit-dev" />
        </el-form-item>
        <el-form-item label="API Secret">
          <el-input v-model="livekit.api_secret" type="password" show-password placeholder="" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="saving.livekit" @click="save('livekit', livekit)">Save LiveKit</el-button>
          <el-tag v-if="savedAt.livekit" type="success" class="ml-2">Saved {{ savedAt.livekit }}</el-tag>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Group 7: Admin Seed (⚠️ dev only) -->
    <el-card>
      <template #header><b>⑧ Admin Seed</b><span class="ml-2 text-xs text-gray-400">section = <code>seed_admin</code> · <span class="text-red-400">dev only — only applies when staff table is empty</span></span></template>
      <el-alert title="⚠️ These credentials are only used on FIRST START when staff table is empty. Changing them here won't reset existing staff passwords." type="warning" show-icon :closable="false" class="mb-3" />
      <el-form :model="seedAdmin" label-width="200px" class="max-w-3xl">
        <el-form-item label="Seed Email">
          <el-input v-model="seedAdmin.email" placeholder="admin@ukteahouse.co.uk" />
        </el-form-item>
        <el-form-item label="Seed Name">
          <el-input v-model="seedAdmin.name" placeholder="Admin" />
        </el-form-item>
        <el-form-item label="Seed Password">
          <el-input v-model="seedAdmin.password" type="password" show-password placeholder="Admin!Tea2026" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="saving.seed_admin" @click="save('seed_admin', seedAdmin)">Save Seed Admin</el-button>
          <el-tag v-if="savedAt.seed_admin" type="success" class="ml-2">Saved {{ savedAt.livekit, seedAdmin }}</el-tag>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Group 8: Runtime Status (read-only) -->
    <el-card>
      <template #header><b>⑨ Runtime Status</b> <span class="text-xs text-gray-400">(read-only)</span></template>
      <div class="grid grid-cols-2 gap-4">
        <div>Config entries in DB: <b>{{ knownKeys.length }}</b></div>
        <div>LiveKit URL: <b>{{ livekit.url || '(not set — using env)' }}</b></div>
        <div>Seed Email: <b>{{ seedAdmin.email || '(not set — using env default)' }}</b></div>
        <div>Current loaded version: <b>{{ server.version || '(not set — using env default)' }}</b></div>
        <div>Service name: <b>{{ app.service_name }}</b></div>
        <div>Magic-link domain: <b>{{ app.domain || '(not set)' }}</b></div>
        <div>Mail Provider: <b>{{ mail.provider || '(not set — default mailgun)' }}</b></div>
        <div>OBS Push: <b>{{ obs.push_domain || '(not set — using default)' }}</b></div>
        <div>OBS Pull: <b>{{ obs.pull_domain || '(not set — using default)' }}</b></div>
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
const knownKeys = ref<string[]>(['app', 'mail', 'obs', 'translate', 'nodes', 'server', 'livekit', 'seed_admin'])

const app      = reactive<Section>({ domain: '', service_name: '', contact_email: '' })
const mail     = reactive<Section>({ provider: '', smtp_host: '', from_addr: '', api_key: '' })
const obs      = reactive<Section>({ push_domain: '', pull_domain: '', rtmp_port: 1935 })
const translate = reactive<Section>({ service_url: '' })
const nodes    = reactive<Section>({ wg_ip_start: '', wg_ip_prefix: '' })
const server   = reactive<Section>({ version: '' })
const livekit  = reactive<Section>({ url: '', api_key: '', api_secret: '' })
const seedAdmin = reactive<Section>({ email: '', name: '', password: '' })

const saving = reactive<Record<string, boolean>>({})
const savedAt = reactive<Record<string, string>>({})

const sections: Record<string, Section> = { app, mail, obs, translate, nodes, server, livekit, seedAdmin }

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

<template>
  <div class="pt-20 py-20 bg-ivory-100">
    <div class="max-w-3xl mx-auto px-6">
      <!-- Section title -->
      <div class="mb-10">
        <div class="flex items-center gap-3 mb-5">
          <span class="h-px w-8 bg-gold/60"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Legal · Privacy Policy</span>
        </div>
        <h1 class="font-serif text-4xl md:text-5xl text-ink-900 mb-4 leading-tight">
          Privacy &amp; <span class="italic">GDPR</span>
        </h1>
        <p class="text-sm text-sand font-serif">Last updated: 18 September 2026 · ICO registration: pending · Company No. 12345678</p>
      </div>

      <!-- Summary — subtle ink panel instead of amber warning box -->
      <div class="bg-ink-900 text-ivory-100 p-6 mb-12" style="border-radius: 2px;">
        <div class="flex items-center gap-2 mb-3">
          <span class="w-1 h-1 rounded-full bg-gold"></span>
          <span class="text-[10px] uppercase tracking-lux text-gold font-sans">Summary</span>
        </div>
        <p class="text-sm text-ivory-100/80 leading-relaxed font-serif">
          We collect only what we need to deliver our service. We never sell your data. Payments are handled by 2Checkout / PayPal — we never see or store your card details. Audit logs of staff actions are kept in a separate, immutable database (GDPR-compliant). You can exercise all GDPR rights via the
          <a href="/gdpr-dsar" class="text-gold underline hover:text-ivory-100">Data Subject Access Request page</a>.
        </p>
      </div>

      <!-- Body -->
      <div class="space-y-10 text-ink-900 leading-relaxed">

        <section>
          <div class="flex items-center gap-3 mb-3">
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">01</span>
            <span class="h-px flex-1 bg-gold/20"></span>
          </div>
          <h2 class="font-serif text-xl text-ink-900 mb-3">Who we are</h2>
          <p class="text-sm text-sand font-serif leading-relaxed">
            UK Tea House Ltd. (company no. 12345678, registered in England &amp; Wales, VAT GB123456789) operates ukteahouse.co.uk. Our Data Protection Officer can be reached at privacy@ukteahouse.co.uk or 12 Savile Row, Mayfair, London W1J 5PA.
          </p>
        </section>

        <section>
          <div class="flex items-center gap-3 mb-3">
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">02</span>
            <span class="h-px flex-1 bg-gold/20"></span>
          </div>
          <h2 class="font-serif text-xl text-ink-900 mb-3">Legal basis for processing</h2>
          <p class="text-sm text-sand font-serif mb-3 leading-relaxed">We process your personal data only under one or more of the following lawful bases:</p>
          <ul class="text-sm text-sand font-serif space-y-2 pl-0 leading-relaxed">
            <li class="grid grid-cols-[auto,1fr] gap-3"><span class="text-gold">—</span><span><strong class="text-ink-900">Contract:</strong> Necessary to fulfil a contract with you (orders, deliveries, customer support).</span></li>
            <li class="grid grid-cols-[auto,1fr] gap-3"><span class="text-gold">—</span><span><strong class="text-ink-900">Legitimate interest:</strong> To run our business securely, improve the website, and communicate service updates.</span></li>
            <li class="grid grid-cols-[auto,1fr] gap-3"><span class="text-gold">—</span><span><strong class="text-ink-900">Consent:</strong> Where you have explicitly agreed (analytics cookies, optional marketing).</span></li>
            <li class="grid grid-cols-[auto,1fr] gap-3"><span class="text-gold">—</span><span><strong class="text-ink-900">Legal obligation:</strong> To comply with tax, accounting, and regulatory requirements.</span></li>
          </ul>
        </section>

        <section>
          <div class="flex items-center gap-3 mb-3">
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">03</span>
            <span class="h-px flex-1 bg-gold/20"></span>
          </div>
          <h2 class="font-serif text-xl text-ink-900 mb-4">What we collect</h2>
          <div class="border border-gold/20" style="border-radius: 2px;">
            <div class="grid grid-cols-3 bg-ink-900 text-ivory-100 text-[10px] uppercase tracking-lux font-sans px-4 py-3">
              <div>Category</div>
              <div>What exactly</div>
              <div>Purpose</div>
            </div>
            <div class="bg-white">
              <div class="grid grid-cols-3 px-4 py-3 text-xs text-sand font-serif border-t border-gold/10">
                <div class="font-medium text-ink-900">Account</div>
                <div>Name, email, hashed password</div>
                <div>Authentication, order history</div>
              </div>
              <div class="grid grid-cols-3 px-4 py-3 text-xs text-sand font-serif border-t border-gold/10">
                <div class="font-medium text-ink-900">Orders</div>
                <div>Shipping/billing address, phone, order items, notes</div>
                <div>Fulfilment, delivery support</div>
              </div>
              <div class="grid grid-cols-3 px-4 py-3 text-xs text-sand font-serif border-t border-gold/10">
                <div class="font-medium text-ink-900">Payment</div>
                <div>Order reference only (no card data)</div>
                <div>Reconciliation with processor</div>
              </div>
              <div class="grid grid-cols-3 px-4 py-3 text-xs text-sand font-serif border-t border-gold/10">
                <div class="font-medium text-ink-900">Technical</div>
                <div>IP address, browser type, cookie IDs</div>
                <div>Fraud prevention, security analytics</div>
              </div>
              <div class="grid grid-cols-3 px-4 py-3 text-xs text-sand font-serif border-t border-gold/10">
                <div class="font-medium text-ink-900">Staff audit</div>
                <div>Action type, subject, timestamp, staff ID (in separate immutable DB)</div>
                <div>GDPR accountability, security</div>
              </div>
            </div>
          </div>
        </section>

        <section>
          <div class="flex items-center gap-3 mb-3">
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">04</span>
            <span class="h-px flex-1 bg-gold/20"></span>
          </div>
          <h2 class="font-serif text-xl text-ink-900 mb-3">What we do NOT collect</h2>
          <ul class="text-sm text-sand font-serif space-y-2 leading-relaxed">
            <li class="grid grid-cols-[auto,1fr] gap-3"><span class="text-gold">—</span><span>Credit card numbers, CVV, or cardholder name — we never touch your payment card. All payments are handled by PCI-DSS certified processors (2Checkout, PayPal).</span></li>
            <li class="grid grid-cols-[auto,1fr] gap-3"><span class="text-gold">—</span><span>Special categories of data (race, religion, health, biometrics) unless you voluntarily provide it in an order note.</span></li>
            <li class="grid grid-cols-[auto,1fr] gap-3"><span class="text-gold">—</span><span>More data than is necessary for each specific purpose.</span></li>
          </ul>
        </section>

        <section>
          <div class="flex items-center gap-3 mb-3">
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">05</span>
            <span class="h-px flex-1 bg-gold/20"></span>
          </div>
          <h2 class="font-serif text-xl text-ink-900 mb-3">Who we share data with</h2>
          <ul class="text-sm text-sand font-serif space-y-2 leading-relaxed">
            <li class="grid grid-cols-[auto,1fr] gap-3"><span class="text-gold">—</span><span>Payment processors (2Checkout, PayPal) — only the minimum required to process payment.</span></li>
            <li class="grid grid-cols-[auto,1fr] gap-3"><span class="text-gold">—</span><span>Shipping couriers — your delivery address and phone number.</span></li>
            <li class="grid grid-cols-[auto,1fr] gap-3"><span class="text-gold">—</span><span>Tax advisors / accountants — anonymised transaction summaries for compliance.</span></li>
            <li class="grid grid-cols-[auto,1fr] gap-3"><span class="text-gold">—</span><span>Law enforcement — only when legally required, with prior notice where permitted.</span></li>
          </ul>
        </section>

        <section>
          <div class="flex items-center gap-3 mb-3">
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">06</span>
            <span class="h-px flex-1 bg-gold/20"></span>
          </div>
          <h2 class="font-serif text-xl text-ink-900 mb-3">International transfers</h2>
          <p class="text-sm text-sand font-serif leading-relaxed">
            Your data is hosted in the European Union (Hetzner EU nodes). Live camera feeds are self-hosted in London. No data is transferred outside the EU except where a payment processor or courier legitimately requires it, in which case standard contractual clauses under GDPR Article 46 apply.
          </p>
        </section>

        <section>
          <div class="flex items-center gap-3 mb-3">
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">07</span>
            <span class="h-px flex-1 bg-gold/20"></span>
          </div>
          <h2 class="font-serif text-xl text-ink-900 mb-4">Data retention</h2>
          <div class="border border-gold/20 bg-white p-5 space-y-2 text-sm text-sand font-serif" style="border-radius: 2px;">
            <div><strong class="text-ink-900">Orders:</strong> 7 years — required by UK VAT law.</div>
            <div><strong class="text-ink-900">Staff audit logs:</strong> 10 years — GDPR accountability + financial compliance.</div>
            <div><strong class="text-ink-900">Account:</strong> Until you delete it, plus 3 years after deletion for tax record purposes.</div>
            <div><strong class="text-ink-900">Cookie consent records:</strong> 5 years.</div>
            <div><strong class="text-ink-900">DSAR request tickets:</strong> 3 years after closure for compliance audit.</div>
          </div>
        </section>

        <section>
          <div class="flex items-center gap-3 mb-3">
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">08</span>
            <span class="h-px flex-1 bg-gold/20"></span>
          </div>
          <h2 class="font-serif text-xl text-ink-900 mb-3">Your rights (GDPR Articles 15–22)</h2>
          <p class="text-sm text-sand font-serif mb-3 leading-relaxed">You have the right to:</p>
          <ul class="text-sm text-sand font-serif space-y-2 leading-relaxed">
            <li class="grid grid-cols-[auto,1fr] gap-3"><span class="text-gold">—</span><span><strong class="text-ink-900">Access</strong> a copy of your data</span></li>
            <li class="grid grid-cols-[auto,1fr] gap-3"><span class="text-gold">—</span><span><strong class="text-ink-900">Rectification</strong> of inaccurate data</span></li>
            <li class="grid grid-cols-[auto,1fr] gap-3"><span class="text-gold">—</span><span><strong class="text-ink-900">Erasure</strong> ("right to be forgotten")</span></li>
            <li class="grid grid-cols-[auto,1fr] gap-3"><span class="text-gold">—</span><span><strong class="text-ink-900">Restrict</strong> processing</span></li>
            <li class="grid grid-cols-[auto,1fr] gap-3"><span class="text-gold">—</span><span><strong class="text-ink-900">Object</strong> to certain processing</span></li>
            <li class="grid grid-cols-[auto,1fr] gap-3"><span class="text-gold">—</span><span><strong class="text-ink-900">Portability</strong> in machine-readable format</span></li>
            <li class="grid grid-cols-[auto,1fr] gap-3"><span class="text-gold">—</span><span><strong class="text-ink-900">Withdraw consent</strong> at any time</span></li>
            <li class="grid grid-cols-[auto,1fr] gap-3"><span class="text-gold">—</span><span><strong class="text-ink-900">Lodge a complaint</strong> with the ICO</span></li>
          </ul>
          <p class="text-sm text-sand font-serif mt-4 leading-relaxed">
            Exercise these rights via our <a href="/gdpr-dsar" class="text-gold underline hover:text-ink-900">Data Subject Access Request form</a> or by emailing privacy@ukteahouse.co.uk. SLA: 30 calendar days from receipt.
          </p>
        </section>

        <section>
          <div class="flex items-center gap-3 mb-3">
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">09</span>
            <span class="h-px flex-1 bg-gold/20"></span>
          </div>
          <h2 class="font-serif text-xl text-ink-900 mb-3">Security</h2>
          <p class="text-sm text-sand font-serif leading-relaxed">
            We use HTTPS everywhere, AES-256 encryption at rest for databases, bcrypt hashing for passwords, and JWT tokens with configurable expiry. Staff audit logs are stored in a separate, append-only database (audit_logs) that cannot be modified or deleted from the business application. Cookie consent submissions are also stored immutably.
          </p>
        </section>

        <section>
          <div class="flex items-center gap-3 mb-3">
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">10</span>
            <span class="h-px flex-1 bg-gold/20"></span>
          </div>
          <h2 class="font-serif text-xl text-ink-900 mb-3">Children</h2>
          <p class="text-sm text-sand font-serif leading-relaxed">
            Our website is not directed at individuals under 18 years of age. We do not knowingly collect data from children. If you are a parent or guardian and believe a child has submitted data, contact us at privacy@ukteahouse.co.uk and we will delete it.
          </p>
        </section>

        <section>
          <div class="flex items-center gap-3 mb-3">
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">11</span>
            <span class="h-px flex-1 bg-gold/20"></span>
          </div>
          <h2 class="font-serif text-xl text-ink-900 mb-3">Changes to this policy</h2>
          <p class="text-sm text-sand font-serif leading-relaxed">
            We may update this policy. Material changes will be notified by email to registered users and highlighted on our website. Continued use of the website after changes constitutes acceptance.
          </p>
        </section>

        <section class="pt-8 border-t border-gold/20">
          <div class="flex items-center gap-3 mb-3">
            <span class="text-[10px] uppercase tracking-lux text-gold font-sans">12</span>
            <span class="h-px flex-1 bg-gold/20"></span>
          </div>
          <h2 class="font-serif text-xl text-ink-900 mb-4">Contact the DPO</h2>
          <div class="text-sm text-sand font-serif space-y-1 leading-relaxed">
            <p>Data Protection Officer</p>
            <p>privacy@ukteahouse.co.uk</p>
            <p>12 Savile Row, Mayfair, London W1J 5PA</p>
            <p class="text-sand/70 pt-2 text-xs">Registered in England &amp; Wales · Company No. 12345678 · VAT GB123456789</p>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

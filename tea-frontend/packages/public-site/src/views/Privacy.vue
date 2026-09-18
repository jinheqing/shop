<template>
  <div class="pt-20 py-20 bg-tea-50">
    <div class="max-w-3xl mx-auto px-6">
      <div class="mb-10">
        <div class="text-xs tracking-[0.25em] uppercase text-tea-500 mb-3">Legal · Privacy Policy</div>
        <h1 class="font-display font-semibold text-4xl md:text-5xl text-tea-900 mb-4">Privacy &amp; <span class="italic">GDPR</span></h1>
        <p class="text-sm text-tea-600">Last updated: 18 September 2026 · ICO registration: pending · Company No. 12345678</p>
      </div>

      <div class="bg-amber-50 border border-amber-200 rounded-2xl p-5 mb-10">
        <div class="text-sm text-amber-900 leading-relaxed">
          <strong>Summary:</strong> We collect only what we need to deliver our service. We never sell your data. Payments are handled by 2Checkout / PayPal — we never see or store your card details. Audit logs of staff actions are kept in a separate, immutable database (GDPR-compliant). You can exercise all GDPR rights via the <a href="/gdpr-dsar" class="underline hover:text-amber-950">Data Subject Access Request page</a>.
        </div>
      </div>

      <div class="space-y-8 text-tea-800 leading-relaxed">
        <section>
          <h2 class="font-display text-xl text-tea-900 mb-2">1. Who we are</h2>
          <p class="text-sm">UK Tea House Ltd. (company no. 12345678, registered in England &amp; Wales, VAT GB123456789) operates ukteahouse.co.uk. Our Data Protection Officer can be reached at privacy@ukteahouse.co.uk or 12 Savile Row, Mayfair, London W1J 5PA.</p>
        </section>

        <section>
          <h2 class="font-display text-xl text-tea-900 mb-2">2. Legal basis for processing</h2>
          <p class="text-sm mb-3">We process your personal data only under one or more of the following lawful bases:</p>
          <ul class="text-sm list-disc list-inside space-y-1 text-tea-700">
            <li><strong>Contract:</strong> Necessary to fulfil a contract with you (orders, deliveries, customer support).</li>
            <li><strong>Legitimate interest:</strong> To run our business securely, improve the website, and communicate service updates.</li>
            <li><strong>Consent:</strong> Where you have explicitly agreed (analytics cookies, optional marketing).</li>
            <li><strong>Legal obligation:</strong> To comply with tax, accounting, and regulatory requirements.</li>
          </ul>
        </section>

        <section>
          <h2 class="font-display text-xl text-tea-900 mb-2">3. What we collect</h2>
          <table class="w-full text-sm border border-tea-200 rounded-xl overflow-hidden">
            <thead class="bg-tea-100 text-tea-900">
              <tr>
                <th class="p-3 text-left">Category</th>
                <th class="p-3 text-left">What exactly</th>
                <th class="p-3 text-left">Purpose</th>
              </tr>
            </thead>
            <tbody class="bg-white text-tea-700">
              <tr class="border-t border-tea-100"><td class="p-3 font-medium">Account</td><td class="p-3">Name, email, hashed password</td><td class="p-3">Authentication, order history</td></tr>
              <tr class="border-t border-tea-100"><td class="p-3 font-medium">Orders</td><td class="p-3">Shipping/billing address, phone, order items, notes</td><td class="p-3">Fulfilment, delivery support</td></tr>
              <tr class="border-t border-tea-100"><td class="p-3 font-medium">Payment</td><td class="p-3">Order reference only (no card data)</td><td class="p-3">Reconciliation with processor</td></tr>
              <tr class="border-t border-tea-100"><td class="p-3 font-medium">Technical</td><td class="p-3">IP address, browser type, cookie IDs</td><td class="p-3">Fraud prevention, security analytics</td></tr>
              <tr class="border-t border-tea-100"><td class="p-3 font-medium">Staff audit</td><td class="p-3">Action type, subject, timestamp, staff ID (in separate immutable DB)</td><td class="p-3">GDPR accountability, security</td></tr>
            </tbody>
          </table>
        </section>

        <section>
          <h2 class="font-display text-xl text-tea-900 mb-2">4. What we do NOT collect</h2>
          <ul class="text-sm list-disc list-inside space-y-1 text-tea-700">
            <li>Credit card numbers, CVV, or cardholder name — we never touch your payment card. All payments are handled by PCI-DSS certified processors (2Checkout, PayPal).</li>
            <li>Special categories of data (race, religion, health, biometrics) unless you voluntarily provide it in an order note.</li>
            <li>More data than is necessary for each specific purpose.</li>
          </ul>
        </section>

        <section>
          <h2 class="font-display text-xl text-tea-900 mb-2">5. Who we share data with</h2>
          <ul class="text-sm list-disc list-inside space-y-1 text-tea-700">
            <li>Payment processors (2Checkout, PayPal) — only the minimum required to process payment.</li>
            <li>Shipping couriers — your delivery address and phone number.</li>
            <li>Tax advisors / accountants — anonymised transaction summaries for compliance.</li>
            <li>Law enforcement — only when legally required, with prior notice where permitted.</li>
          </ul>
        </section>

        <section>
          <h2 class="font-display text-xl text-tea-900 mb-2">6. International transfers</h2>
          <p class="text-sm">Your data is hosted in the European Union (Hetzner EU nodes). Live camera feeds are self-hosted in London. No data is transferred outside the EU except where a payment processor or courier legitimately requires it, in which case standard contractual clauses under GDPR Article 46 apply.</p>
        </section>

        <section>
          <h2 class="font-display text-xl text-tea-900 mb-2">7. Data retention</h2>
          <div class="rounded-xl border border-tea-100 bg-tea-50/50 p-5 text-sm text-tea-700 space-y-1">
            <div><strong>Orders:</strong> 7 years — required by UK VAT law.</div>
            <div><strong>Staff audit logs:</strong> 10 years — GDPR accountability + financial compliance.</div>
            <div><strong>Account:</strong> Until you delete it, plus 3 years after deletion for tax record purposes.</div>
            <div><strong>Cookie consent records:</strong> 5 years.</div>
            <div><strong>DSAR request tickets:</strong> 3 years after closure for compliance audit.</div>
          </div>
        </section>

        <section>
          <h2 class="font-display text-xl text-tea-900 mb-2">8. Your rights (GDPR Articles 15–22)</h2>
          <p class="text-sm mb-3">You have the right to:</p>
          <ul class="text-sm list-disc list-inside space-y-1 text-tea-700">
            <li><strong>Access</strong> a copy of your data</li>
            <li><strong>Rectification</strong> of inaccurate data</li>
            <li><strong>Erasure</strong> ("right to be forgotten")</li>
            <li><strong>Restrict</strong> processing</li>
            <li><strong>Object</strong> to certain processing</li>
            <li><strong>Portability</strong> in machine-readable format</li>
            <li><strong>Withdraw consent</strong> at any time</li>
            <li><strong>Lodge a complaint</strong> with the ICO</li>
          </ul>
          <p class="text-sm mt-3">Exercise these rights via our <a href="/gdpr-dsar" class="text-tea-800 underline hover:text-tea-950">Data Subject Access Request form</a> or by emailing privacy@ukteahouse.co.uk. SLA: 30 calendar days from receipt.</p>
        </section>

        <section>
          <h2 class="font-display text-xl text-tea-900 mb-2">9. Security</h2>
          <p class="text-sm">We use HTTPS everywhere, AES-256 encryption at rest for databases, bcrypt hashing for passwords, and JWT tokens with configurable expiry. Staff audit logs are stored in a separate, append-only database (audit_logs) that cannot be modified or deleted from the business application. Cookie consent submissions are also stored immutably.</p>
        </section>

        <section>
          <h2 class="font-display text-xl text-tea-900 mb-2">10. Children</h2>
          <p class="text-sm">Our website is not directed at individuals under 18 years of age. We do not knowingly collect data from children. If you are a parent or guardian and believe a child has submitted data, contact us at privacy@ukteahouse.co.uk and we will delete it.</p>
        </section>

        <section>
          <h2 class="font-display text-xl text-tea-900 mb-2">11. Changes to this policy</h2>
          <p class="text-sm">We may update this policy. Material changes will be notified by email to registered users and highlighted on our website. Continued use of the website after changes constitutes acceptance.</p>
        </section>

        <section class="pt-4 border-t border-tea-200">
          <h2 class="font-display text-xl text-tea-900 mb-3">12. Contact the DPO</h2>
          <p class="text-sm mb-1">Data Protection Officer</p>
          <p class="text-sm mb-1">privacy@ukteahouse.co.uk</p>
          <p class="text-sm mb-1">12 Savile Row, Mayfair, London W1J 5PA</p>
          <p class="text-sm text-tea-500">Registered in England &amp; Wales · Company No. 12345678 · VAT GB123456789</p>
        </section>
      </div>
    </div>
  </div>
</template>

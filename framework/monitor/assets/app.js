(function () {
  'use strict';

  // ---------- mini-chart ---------------------------------------------------
  // Time series sobre <canvas>. Suporta múltiplas séries com escala Y
  // automática, eixo X temporal, grade leve, hover com tooltip e legenda.
  function MiniChart(canvas, opts) {
    this.canvas = canvas;
    this.opts = opts || {};
    this.series = (opts.series || []).map(function (s) {
      return { label: s.label, color: s.color, axis: s.axis || 'left', data: [] };
    });
    this.dpr = window.devicePixelRatio || 1;
    this.padding = { top: 10, right: 50, bottom: 22, left: 50 };
    this.maxPoints = opts.maxPoints || 600;
    this.yFormat = opts.yFormat || function (v) { return v.toFixed(0); };
    this.yFormatRight = opts.yFormatRight || this.yFormat;
    this.hover = null;
    this.resize();
    var self = this;
    window.addEventListener('resize', function () { self.resize(); self.draw(); });
    canvas.addEventListener('mousemove', function (e) { self.onMove(e); });
    canvas.addEventListener('mouseleave', function () { self.hover = null; self.draw(); });
  }

  MiniChart.prototype.resize = function () {
    var rect = this.canvas.getBoundingClientRect();
    this.w = rect.width;
    this.h = rect.height;
    this.canvas.width = this.w * this.dpr;
    this.canvas.height = this.h * this.dpr;
    this.canvas.style.height = this.h + 'px';
  };

  MiniChart.prototype.push = function (t, values) {
    for (var i = 0; i < this.series.length; i++) {
      var v = values[i];
      if (v === undefined || v === null) continue;
      this.series[i].data.push({ t: t, v: v });
      if (this.series[i].data.length > this.maxPoints) {
        this.series[i].data.shift();
      }
    }
  };

  MiniChart.prototype.setAll = function (samples, extractors) {
    for (var i = 0; i < this.series.length; i++) this.series[i].data = [];
    for (var j = 0; j < samples.length; j++) {
      var s = samples[j];
      for (var k = 0; k < extractors.length; k++) {
        var v = extractors[k](s);
        if (v === undefined || v === null) continue;
        this.series[k].data.push({ t: s.t, v: v });
      }
    }
    while (this.series[0] && this.series[0].data.length > this.maxPoints) {
      for (var n = 0; n < this.series.length; n++) this.series[n].data.shift();
    }
  };

  MiniChart.prototype.bounds = function (axis) {
    var min = Infinity, max = -Infinity;
    for (var i = 0; i < this.series.length; i++) {
      if (this.series[i].axis !== axis) continue;
      for (var j = 0; j < this.series[i].data.length; j++) {
        var v = this.series[i].data[j].v;
        if (v < min) min = v;
        if (v > max) max = v;
      }
    }
    if (min === Infinity) { min = 0; max = 1; }
    if (min === max) { max = min + 1; }
    var pad = (max - min) * 0.1;
    return { min: Math.max(0, min - pad), max: max + pad };
  };

  MiniChart.prototype.tBounds = function () {
    var min = Infinity, max = -Infinity;
    for (var i = 0; i < this.series.length; i++) {
      var d = this.series[i].data;
      if (!d.length) continue;
      if (d[0].t < min) min = d[0].t;
      if (d[d.length - 1].t > max) max = d[d.length - 1].t;
    }
    if (min === Infinity) { var now = Date.now(); min = now - 60000; max = now; }
    if (min === max) max = min + 1;
    return { min: min, max: max };
  };

  MiniChart.prototype.xy = function (t, v, axis) {
    var tb = this._tb;
    var yb = axis === 'right' ? this._ybR : this._ybL;
    var x = this.padding.left + (t - tb.min) / (tb.max - tb.min) * (this.w - this.padding.left - this.padding.right);
    var y = this.padding.top + (1 - (v - yb.min) / (yb.max - yb.min)) * (this.h - this.padding.top - this.padding.bottom);
    return { x: x, y: y };
  };

  MiniChart.prototype.draw = function () {
    var ctx = this.canvas.getContext('2d');
    ctx.setTransform(this.dpr, 0, 0, this.dpr, 0, 0);
    ctx.clearRect(0, 0, this.w, this.h);

    this._tb = this.tBounds();
    this._ybL = this.bounds('left');
    this._ybR = this.bounds('right');

    var hasRight = this.series.some(function (s) { return s.axis === 'right'; });

    // grade horizontal e ticks
    ctx.strokeStyle = '#2a313c';
    ctx.fillStyle = '#8b949e';
    ctx.lineWidth = 1;
    ctx.font = '10px ui-monospace, Menlo, Consolas, monospace';
    ctx.textBaseline = 'middle';
    var gridSteps = 4;
    for (var g = 0; g <= gridSteps; g++) {
      var ratio = g / gridSteps;
      var y = this.padding.top + ratio * (this.h - this.padding.top - this.padding.bottom);
      ctx.beginPath();
      ctx.moveTo(this.padding.left, y);
      ctx.lineTo(this.w - this.padding.right, y);
      ctx.stroke();
      var leftV = this._ybL.max - (this._ybL.max - this._ybL.min) * ratio;
      ctx.textAlign = 'right';
      ctx.fillText(this.yFormat(leftV), this.padding.left - 4, y);
      if (hasRight) {
        var rightV = this._ybR.max - (this._ybR.max - this._ybR.min) * ratio;
        ctx.textAlign = 'left';
        ctx.fillText(this.yFormatRight(rightV), this.w - this.padding.right + 4, y);
      }
    }

    // ticks de tempo
    ctx.textAlign = 'center';
    ctx.textBaseline = 'top';
    var tSteps = 4;
    for (var t = 0; t <= tSteps; t++) {
      var rt = t / tSteps;
      var ts = this._tb.min + (this._tb.max - this._tb.min) * rt;
      var x = this.padding.left + rt * (this.w - this.padding.left - this.padding.right);
      var d = new Date(ts);
      var label = pad2(d.getHours()) + ':' + pad2(d.getMinutes()) + ':' + pad2(d.getSeconds());
      ctx.fillText(label, x, this.h - this.padding.bottom + 4);
    }

    // séries
    for (var i = 0; i < this.series.length; i++) {
      var s = this.series[i];
      if (!s.data.length) continue;
      ctx.strokeStyle = s.color;
      ctx.lineWidth = 1.5;
      ctx.beginPath();
      for (var j = 0; j < s.data.length; j++) {
        var p = this.xy(s.data[j].t, s.data[j].v, s.axis);
        if (j === 0) ctx.moveTo(p.x, p.y); else ctx.lineTo(p.x, p.y);
      }
      ctx.stroke();
      // último ponto
      var last = s.data[s.data.length - 1];
      var lp = this.xy(last.t, last.v, s.axis);
      ctx.fillStyle = s.color;
      ctx.beginPath();
      ctx.arc(lp.x, lp.y, 2.5, 0, Math.PI * 2);
      ctx.fill();
    }

    // legenda
    var lx = this.padding.left + 4;
    var ly = this.padding.top + 4;
    ctx.font = '11px -apple-system, "Segoe UI", system-ui, sans-serif';
    ctx.textBaseline = 'top';
    for (var k = 0; k < this.series.length; k++) {
      var s2 = this.series[k];
      ctx.fillStyle = s2.color;
      ctx.fillRect(lx, ly + 3, 10, 2);
      ctx.fillStyle = '#e6edf3';
      ctx.textAlign = 'left';
      var lv = s2.data.length ? s2.data[s2.data.length - 1].v : 0;
      var fmt = s2.axis === 'right' ? this.yFormatRight : this.yFormat;
      var label = s2.label + ' ' + fmt(lv);
      ctx.fillText(label, lx + 14, ly);
      lx += ctx.measureText(label).width + 28;
    }

    // hover crosshair + tooltip
    if (this.hover) {
      ctx.strokeStyle = '#8b949e';
      ctx.setLineDash([3, 3]);
      ctx.beginPath();
      ctx.moveTo(this.hover.x, this.padding.top);
      ctx.lineTo(this.hover.x, this.h - this.padding.bottom);
      ctx.stroke();
      ctx.setLineDash([]);
      this.drawTooltip(ctx, this.hover);
    }
  };

  MiniChart.prototype.onMove = function (e) {
    var rect = this.canvas.getBoundingClientRect();
    var x = e.clientX - rect.left;
    if (x < this.padding.left || x > this.w - this.padding.right) {
      this.hover = null; this.draw(); return;
    }
    this.hover = { x: x };
    this.draw();
  };

  MiniChart.prototype.drawTooltip = function (ctx, hover) {
    if (!this.series.length) return;
    var tb = this._tb;
    var ratio = (hover.x - this.padding.left) / (this.w - this.padding.left - this.padding.right);
    var t = tb.min + ratio * (tb.max - tb.min);
    // pega o ponto mais próximo de t na primeira série
    var lines = [];
    var ts = null;
    for (var i = 0; i < this.series.length; i++) {
      var s = this.series[i];
      var nearest = nearestPoint(s.data, t);
      if (!nearest) continue;
      ts = nearest.t;
      var fmt = s.axis === 'right' ? this.yFormatRight : this.yFormat;
      lines.push({ color: s.color, text: s.label + ': ' + fmt(nearest.v) });
    }
    if (!lines.length) return;
    if (ts) {
      var d = new Date(ts);
      lines.unshift({ color: '#8b949e', text: pad2(d.getHours()) + ':' + pad2(d.getMinutes()) + ':' + pad2(d.getSeconds()) });
    }
    ctx.font = '11px -apple-system, "Segoe UI", system-ui, sans-serif';
    var w = 0;
    for (var j = 0; j < lines.length; j++) w = Math.max(w, ctx.measureText(lines[j].text).width);
    w += 22;
    var h = 6 + lines.length * 14;
    var tx = hover.x + 10;
    if (tx + w > this.w - this.padding.right) tx = hover.x - w - 10;
    var ty = this.padding.top + 4;
    ctx.fillStyle = 'rgba(13,17,23,0.95)';
    ctx.strokeStyle = '#2a313c';
    roundRect(ctx, tx, ty, w, h, 4);
    ctx.fill();
    ctx.stroke();
    for (var k = 0; k < lines.length; k++) {
      ctx.fillStyle = lines[k].color;
      ctx.fillRect(tx + 6, ty + 5 + k * 14 + 4, 8, 2);
      ctx.fillStyle = '#e6edf3';
      ctx.textBaseline = 'top';
      ctx.textAlign = 'left';
      ctx.fillText(lines[k].text, tx + 18, ty + 5 + k * 14);
    }
  };

  function nearestPoint(data, t) {
    if (!data.length) return null;
    var best = data[0];
    var bestDt = Math.abs(data[0].t - t);
    for (var i = 1; i < data.length; i++) {
      var dt = Math.abs(data[i].t - t);
      if (dt < bestDt) { best = data[i]; bestDt = dt; }
    }
    return best;
  }

  function roundRect(ctx, x, y, w, h, r) {
    ctx.beginPath();
    ctx.moveTo(x + r, y);
    ctx.arcTo(x + w, y, x + w, y + h, r);
    ctx.arcTo(x + w, y + h, x, y + h, r);
    ctx.arcTo(x, y + h, x, y, r);
    ctx.arcTo(x, y, x + w, y, r);
    ctx.closePath();
  }

  function pad2(n) { return n < 10 ? '0' + n : '' + n; }

  // ---------- formatters ---------------------------------------------------
  function fmtBytes(n) {
    if (n === undefined || n === null) return '—';
    if (n < 1024) return n + ' B';
    if (n < 1024 * 1024) return (n / 1024).toFixed(1) + ' KB';
    if (n < 1024 * 1024 * 1024) return (n / 1024 / 1024).toFixed(1) + ' MB';
    return (n / 1024 / 1024 / 1024).toFixed(2) + ' GB';
  }
  function fmtMB(n) { return (n / 1024 / 1024).toFixed(1); }
  function fmtUptime(sec) {
    if (!sec) return '—';
    var d = Math.floor(sec / 86400);
    var h = Math.floor((sec % 86400) / 3600);
    var m = Math.floor((sec % 3600) / 60);
    var s = Math.floor(sec % 60);
    var out = '';
    if (d) out += d + 'd ';
    if (d || h) out += h + 'h ';
    if (d || h || m) out += m + 'm ';
    out += s + 's';
    return out;
  }

  // ---------- charts -------------------------------------------------------
  var chartHeap = new MiniChart(document.getElementById('chartHeap'), {
    series: [
      { label: 'Alloc', color: '#58a6ff' },
      { label: 'Inuse', color: '#3fb950' },
      { label: 'Sys',   color: '#d29922' }
    ],
    yFormat: function (v) { return v.toFixed(1); }
  });
  var chartGoroutines = new MiniChart(document.getElementById('chartGoroutines'), {
    series: [
      { label: 'Goroutines', color: '#58a6ff' },
      { label: 'Threads',    color: '#bc8cff', axis: 'right' }
    ]
  });
  var chartGC = new MiniChart(document.getElementById('chartGC'), {
    series: [
      { label: 'Pause ms', color: '#f85149' },
      { label: '%CPU GC',  color: '#d29922', axis: 'right' }
    ],
    yFormat: function (v) { return v.toFixed(1); },
    yFormatRight: function (v) { return v.toFixed(2) + '%'; }
  });
  var chartAlloc = new MiniChart(document.getElementById('chartAlloc'), {
    series: [
      { label: 'Aloc MB/s', color: '#3fb950' }
    ],
    yFormat: function (v) { return v.toFixed(2); }
  });

  function ingest(s) {
    chartHeap.push(s.t, [fmtMB(s.heapAlloc), fmtMB(s.heapInuse), fmtMB(s.heapSys)]);
    chartGoroutines.push(s.t, [s.goroutines, s.threads]);
    chartGC.push(s.t, [s.gcPauseMs, s.gcCpuPercent]);
    chartAlloc.push(s.t, [s.allocRate / 1024 / 1024]);
  }

  function bulkLoad(samples) {
    chartHeap.setAll(samples, [
      function (s) { return fmtMB(s.heapAlloc); },
      function (s) { return fmtMB(s.heapInuse); },
      function (s) { return fmtMB(s.heapSys); }
    ]);
    chartGoroutines.setAll(samples, [
      function (s) { return s.goroutines; },
      function (s) { return s.threads; }
    ]);
    chartGC.setAll(samples, [
      function (s) { return s.gcPauseMs; },
      function (s) { return s.gcCpuPercent; }
    ]);
    chartAlloc.setAll(samples, [
      function (s) { return s.allocRate / 1024 / 1024; }
    ]);
    redrawAll();
  }
  function redrawAll() { chartHeap.draw(); chartGoroutines.draw(); chartGC.draw(); chartAlloc.draw(); }

  // ---------- snapshot/cards -----------------------------------------------
  function applySnapshot(snap) {
    document.getElementById('appName').textContent = snap.appName || 'Monitor';
    document.title = (snap.appName || 'Monitor') + ' — Monitor';
    document.getElementById('userName').textContent = ''; // populado via /api/me futuramente
    document.getElementById('cardGoroutines').textContent = snap.goroutines;
    document.getElementById('cardHeap').textContent = fmtBytes(snap.heapInuse);
    document.getElementById('cardSys').textContent = fmtBytes(snap.sys);
    document.getElementById('cardThreads').textContent = snap.threads;
    document.getElementById('cardGC').textContent = snap.numGC;
    document.getElementById('cardUptime').textContent = fmtUptime(snap.uptime);
    document.getElementById('infoApp').textContent = snap.appName || '';
    document.getElementById('infoHost').textContent = snap.hostname || '';
    document.getElementById('infoPid').textContent = snap.pid;
    document.getElementById('infoGo').textContent = snap.goVersion;
    document.getElementById('infoCpu').textContent = snap.numCPU;
    document.getElementById('infoMaxProc').textContent = snap.gomaxprocs;
    if (snap.alerts && snap.alerts.length) {
      for (var i = 0; i < snap.alerts.length; i++) prependAlert(snap.alerts[i]);
      bumpAlertBadge(snap.alerts.length);
      toast(snap.alerts[0].message, 'err');
    }
  }

  // ---------- SSE ---------------------------------------------------------
  var es = null;
  var status = document.getElementById('connStatus');
  function connect() {
    status.dataset.state = 'connecting';
    status.textContent = 'conectando…';
    es = new EventSource('/api/stream');
    es.addEventListener('snapshot', function (ev) {
      status.dataset.state = 'ok';
      status.textContent = 'live';
      try {
        var snap = JSON.parse(ev.data);
        applySnapshot(snap);
        ingest(snap);
        redrawAll();
      } catch (e) { console.error(e); }
    });
    es.onerror = function () {
      status.dataset.state = 'lost';
      status.textContent = 'reconectando…';
      es.close();
      setTimeout(connect, 2000);
    };
  }

  // ---------- alertas -----------------------------------------------------
  var alertBadgeCount = 0;
  function bumpAlertBadge(n) {
    alertBadgeCount += n;
    var el = document.getElementById('alertBadge');
    el.textContent = alertBadgeCount;
    el.hidden = alertBadgeCount === 0;
  }
  function prependAlert(a) {
    var list = document.getElementById('alertList');
    var empty = list.querySelector('.empty');
    if (empty) empty.remove();
    var div = document.createElement('div');
    div.className = 'alert-item severity-' + (a.severity || 'warn');
    var d = new Date(a.t);
    div.innerHTML = '<span class="when">' + pad2(d.getHours()) + ':' + pad2(d.getMinutes()) + ':' + pad2(d.getSeconds()) + '</span>' +
      '<span class="sev">' + (a.severity || 'warn') + '</span>' +
      '<span class="msg"></span>';
    div.querySelector('.msg').textContent = a.message;
    list.insertBefore(div, list.firstChild);
    while (list.children.length > 200) list.removeChild(list.lastChild);
  }
  function loadAlerts() {
    fetch('/api/alerts').then(function (r) { return r.json(); }).then(function (list) {
      if (!list || !list.length) return;
      var ul = document.getElementById('alertList');
      ul.innerHTML = '';
      for (var i = list.length - 1; i >= 0; i--) prependAlert(list[i]);
    });
  }

  // ---------- regras ------------------------------------------------------
  function ruleRow(r) {
    var tr = document.createElement('tr');
    tr.innerHTML =
      '<td><input type="checkbox" data-k="ativo"></td>' +
      '<td><input type="text" data-k="name"></td>' +
      '<td><input type="text" data-k="metric"></td>' +
      '<td><select data-k="op"><option>&gt;</option><option>&gt;=</option><option>&lt;</option><option>&lt;=</option><option>growth</option></select></td>' +
      '<td><input type="number" step="any" data-k="threshold"></td>' +
      '<td><input type="number" data-k="windowSec"></td>' +
      '<td><select data-k="severity"><option>info</option><option>warn</option><option>critical</option></select></td>' +
      '<td><button class="btn danger" type="button">remover</button></td>';
    tr.querySelector('[data-k="ativo"]').checked = !!r.ativo;
    tr.querySelector('[data-k="name"]').value = r.name || '';
    tr.querySelector('[data-k="metric"]').value = r.metric || '';
    tr.querySelector('[data-k="op"]').value = r.op || '>';
    tr.querySelector('[data-k="threshold"]').value = r.threshold || 0;
    tr.querySelector('[data-k="windowSec"]').value = r.windowSec || 0;
    tr.querySelector('[data-k="severity"]').value = r.severity || 'warn';
    tr.querySelector('button').addEventListener('click', function () { tr.remove(); });
    return tr;
  }
  function loadRules() {
    fetch('/api/rules').then(function (r) { return r.json(); }).then(function (rules) {
      var body = document.getElementById('rulesBody');
      body.innerHTML = '';
      (rules || []).forEach(function (r) { body.appendChild(ruleRow(r)); });
    });
  }
  function saveRules() {
    var rows = document.querySelectorAll('#rulesBody tr');
    var rules = [];
    rows.forEach(function (tr) {
      rules.push({
        ativo: tr.querySelector('[data-k="ativo"]').checked,
        name: tr.querySelector('[data-k="name"]').value,
        metric: tr.querySelector('[data-k="metric"]').value,
        op: tr.querySelector('[data-k="op"]').value,
        threshold: parseFloat(tr.querySelector('[data-k="threshold"]').value) || 0,
        windowSec: parseInt(tr.querySelector('[data-k="windowSec"]').value, 10) || 0,
        severity: tr.querySelector('[data-k="severity"]').value
      });
    });
    fetch('/api/rules', { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(rules) })
      .then(function (r) { if (!r.ok) return r.text().then(function (t) { throw new Error(t); }); return r.json(); })
      .then(function () { toast('regras salvas', 'ok'); })
      .catch(function (e) { toast(e.message, 'err'); });
  }

  // ---------- abas / ações ------------------------------------------------
  document.querySelectorAll('.tab').forEach(function (t) {
    t.addEventListener('click', function () {
      document.querySelectorAll('.tab').forEach(function (x) { x.classList.remove('active'); });
      t.classList.add('active');
      var name = t.dataset.tab;
      document.querySelectorAll('.panel').forEach(function (p) { p.hidden = p.dataset.panel !== name; });
      if (name === 'alerts') loadAlerts();
      if (name === 'rules') loadRules();
    });
  });

  document.getElementById('btnGC').addEventListener('click', function () {
    fetch('/api/gc', { method: 'POST' }).then(function () { toast('GC executado', 'ok'); });
  });
  document.getElementById('btnFree').addEventListener('click', function () {
    fetch('/api/freeosmemory', { method: 'POST' }).then(function () { toast('FreeOSMemory executado', 'ok'); });
  });
  document.getElementById('btnDump').addEventListener('click', function () {
    window.open('/api/goroutines', '_blank');
  });
  document.getElementById('addRule').addEventListener('click', function () {
    document.getElementById('rulesBody').appendChild(ruleRow({ ativo: true, op: '>', severity: 'warn' }));
  });
  document.getElementById('saveRules').addEventListener('click', saveRules);

  document.getElementById('passForm').addEventListener('submit', function (ev) {
    ev.preventDefault();
    var f = ev.target;
    var msg = document.getElementById('passMsg');
    var body = {
      user: f.querySelector('[name="user"]').value || '',
      current: f.querySelector('[name="current"]').value,
      new: f.querySelector('[name="new"]').value
    };
    fetch('/api/password', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body) })
      .then(function (r) {
        if (r.ok) { msg.className = 'msg ok'; msg.textContent = 'senha alterada'; f.reset(); return; }
        return r.text().then(function (t) { msg.className = 'msg err'; msg.textContent = t; });
      })
      .catch(function (e) { msg.className = 'msg err'; msg.textContent = e.message; });
  });

  // ---------- toast -------------------------------------------------------
  var toastTimer;
  function toast(text, kind) {
    var el = document.getElementById('toast');
    el.className = 'toast ' + (kind || '');
    el.textContent = text;
    el.hidden = false;
    clearTimeout(toastTimer);
    toastTimer = setTimeout(function () { el.hidden = true; }, 4000);
  }

  // ---------- bootstrap ---------------------------------------------------
  fetch('/api/history').then(function (r) { return r.json(); }).then(function (h) {
    if (h && h.length) bulkLoad(h);
    connect();
  }).catch(function () { connect(); });
})();

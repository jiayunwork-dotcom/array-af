'use strict';

function $(id) { return document.getElementById(id); }

function showError(msg) {
  const box = $('error-box');
  box.textContent = msg;
  box.classList.remove('hidden');
}

function clearError() {
  $('error-box').classList.add('hidden');
  $('error-box').textContent = '';
}

function readForm() {
  return {
    N: parseInt($('N').value, 10),
    d: parseFloat($('d').value),
    lambda: parseFloat($('lambda').value),
    beta_deg: parseFloat($('beta').value),
    element: $('element').value,
    theta_steps: parseInt($('theta_steps').value, 10)
  };
}

function fillForm(payload) {
  $('N').value = payload.N;
  $('d').value = payload.d;
  $('lambda').value = payload.lambda;
  $('beta').value = (payload.beta_deg !== undefined ? payload.beta_deg : 0);
  $('element').value = payload.element || 'iso';
}

async function postJSON(url, body) {
  const resp = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  });
  const text = await resp.text();
  let data = null;
  try { data = JSON.parse(text); } catch (e) { data = { raw: text }; }
  if (!resp.ok) {
    const msg = data && data.error ? data.error : ('HTTP ' + resp.status + ': ' + text);
    throw new Error(msg);
  }
  return data;
}

function fmt(v, digits) {
  if (v === null || v === undefined || Number.isNaN(v)) return '—';
  return v.toFixed(digits === undefined ? 2 : digits);
}

function renderResult(d) {
  const m = d.mainlobe;
  const h = d.hpbw;
  const g = d.grating;
  const p = d.params;
  const rows = [
    ['主瓣角', m.visible ? fmt(m.angle_deg) + '°' : '不可见（虚空间）'],
    ['HPBW', h.measurable ? fmt(h.width_deg) + '°' : '不可测'],
    ['HPBW 区间', h.measurable ? fmt(h.left_deg) + '° … ' + fmt(h.right_deg) + '°' : '—'],
    ['栅瓣', g.present ? '有 (' + g.lobes.length + ' 个)' : '无'],
    ['AF 峰值', fmt(d.af_peak) + '（匹配 N=' + p.N + '：' + (d.af_peak_matches_n ? '是' : '否') + '）'],
    ['方向性近似', d.directivity.valid ? 'D≈' + fmt(d.directivity.approx, 0) : '不适用（' + d.directivity.reason + '）'],
    ['第一零点', d.nulls.left_valid || d.nulls.right_valid
      ? fmt(d.nulls.left_deg) + '° / ' + fmt(d.nulls.right_deg) + '°'
      : '可见区无'],
    ['d/λ', fmt(p.d / p.lambda, 3)],
    ['端射 β(−kd)', fmt(-(2 * Math.PI / p.lambda) * p.d * 180 / Math.PI) + '°']
  ];
  $('result-grid').innerHTML = rows.map(r =>
    '<div class="rk">' + r[0] + '</div><div class="rv">' + r[1] + '</div>'
  ).join('');
  renderPlot(d);
}

function renderPlot(d) {
  const size = 400;
  const cx = size / 2, cy = size / 2;
  const rMax = size / 2 - 12;
  let svg = '<svg width="' + size + '" height="' + size + '" viewBox="0 0 ' + size + ' ' + size + '">';
  // 坐标参考圆环（r=0.25/0.5/0.75/1 的归一化刻度，非数据）。
  [0.25, 0.5, 0.75, 1].forEach(function (f) {
    svg += '<circle cx="' + cx + '" cy="' + cy + '" r="' + (f * rMax) +
      '" class="grid-ring" stroke="#ccc" fill="none" stroke-dasharray="2 3"/>';
  });
  svg += '<line x1="' + cx + '" y1="0" x2="' + cx + '" y2="' + size +
    '" class="axis" stroke="#ddd"/><line x1="0" y1="' + cy + '" x2="' + size + '" y2="' + cy +
    '" class="axis" stroke="#ddd"/>';
  const pts = d.points || [];
  let path = '';
  pts.forEach(function (p) {
    const x = cx + p.x * rMax;
    const y = cy - p.y * rMax;
    path += (path ? ' L' : 'M') + x.toFixed(2) + ',' + y.toFixed(2);
  });
  svg += '<path d="' + path + '" class="af-curve" fill="none" stroke="#1a63c8" stroke-width="1.5"/>';
  // 主瓣方向标记。
  if (d.mainlobe.visible) {
    const ang = d.mainlobe.angle_deg * Math.PI / 180;
    const x = cx + 0.95 * rMax * Math.sin(ang);
    const y = cy - 0.95 * rMax * Math.cos(ang);
    svg += '<circle cx="' + x.toFixed(2) + '" cy="' + y.toFixed(2) + '" r="4" class="ml-marker" fill="#c0392b"/>';
  }
  svg += '<text x="' + (cx + 6) + '" y="' + (cy - 6) + '" class="lbl">θ=90° 侧射</text></svg>';
  $('plot').innerHTML = svg;
}

function renderScan(d) {
  let s = '<p>β 从 ' + fmt(d.start_beta_deg) + '° 扫到 ' + fmt(d.end_beta_deg) + '°（' + d.steps + ' 段）。</p>';
  const sm = d.summary;
  s += '<ul><li>主瓣侧射 → 端射：' + (sm.mainlobe_moves_toward_endfire ? '是' : '否') +
    '（' + fmt(sm.mainlobe_start_deg) + '° → ' + fmt(sm.mainlobe_end_deg) + '°）</li>' +
    '<li>HPBW 变宽：' + (sm.hpbw_widened ? '是' : '否') + '</li>' +
    '<li>栅瓣出现：' + (sm.grating_appears ? '是' : '否') + '</li></ul>';
  $('scan-summary').innerHTML = s;
  let t = '<table><thead><tr><th>β(°)</th><th>主瓣角(°)</th><th>可见</th><th>HPBW(°)</th><th>栅瓣</th></tr></thead><tbody>';
  d.rows.forEach(function (r) {
    t += '<tr><td>' + fmt(r.beta_deg) + '</td><td>' + fmt(r.mainlobe_deg) + '</td><td>' +
      (r.mainlobe_visible ? '是' : '否') + '</td><td>' + fmt(r.hpbw_deg) + '</td><td>' +
      (r.has_grating ? '有' : '无') + '</td></tr>';
  });
  t += '</tbody></table>';
  $('scan-table-wrap').innerHTML = t;
}

$('af-form').addEventListener('submit', async function (ev) {
  ev.preventDefault();
  clearError();
  try {
    const data = await postJSON('/api/af', readForm());
    renderResult(data);
  } catch (e) {
    showError('计算失败：' + e.message);
  }
});

$('load-example').addEventListener('click', async function () {
  clearError();
  try {
    const resp = await fetch('/api/examples');
    const data = await resp.json();
    if (!resp.ok) { throw new Error(data.error || ('HTTP ' + resp.status)); }
    const payload = data.payloads['broadside-8'];
    if (!payload) { throw new Error('示例 broadside-8 不存在'); }
    fillForm(payload);
  } catch (e) {
    showError('加载示例失败：' + e.message);
  }
});

$('scan-btn').addEventListener('click', async function () {
  clearError();
  try {
    const body = readForm();
    delete body.beta_deg;
    delete body.theta_steps;
    const data = await postJSON('/api/scan', body);
    renderScan(data);
  } catch (e) {
    showError('扫描失败：' + e.message);
  }
});

// 首次进入即用示例算一遍，方便直接看到结果。
(async function () {
  try {
    const data = await postJSON('/api/af', readForm());
    renderResult(data);
  } catch (e) {
    showError('初始计算失败：' + e.message);
  }
})();

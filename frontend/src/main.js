import { OpenFileDialog, SaveFileDialog } from '../wailsjs/go/gui/App';

/* Permissions definition */
const PERMISSIONS = [
    { id: 'perm-print', flag: 'print', label: 'Print', defaultOn: true },
    { id: 'perm-print-low', flag: 'print-low', label: 'Print (low quality)', defaultOn: false },
    { id: 'perm-copy', flag: 'copy', label: 'Copy', defaultOn: true },
    { id: 'perm-extract', flag: 'extract', label: 'Extract', defaultOn: false },
    { id: 'perm-modify', flag: 'modify', label: 'Modify', defaultOn: false },
    { id: 'perm-annotate', flag: 'annotate', label: 'Annotate', defaultOn: true },
    { id: 'perm-fill', flag: 'fill', label: 'Fill forms', defaultOn: false }, 
    { id: 'perm-assemble', flag: 'assemble', label: 'Assemble', defaultOn: false },
];

/* Build permissions grid */
const grid = document.getElementById('permissions-grid');
PERMISSIONS.forEach(p => {
    const item = document.createElement('label');
    item.id = p.id;
    item.className = 'perm-item' + (p.defaultOn ? ' checked' : '');
    item.innerHTML = `
        <input type="checkbox" name="permissions" value="${p.flag}" ${p.defaultOn ? 'checked' : ''} />
        <span class="perm-check">${p.defaultOn ? '<i class="bi bi-check"></i>' : ''}</span>
        <span class="perm-name">${p.label}</span>
      `;
    item.addEventListener('click', () => togglePerm(item));
    grid.appendChild(item);
});

function togglePerm(item) {
    const cb = item.querySelector('input[type=checkbox]');
    const chk = item.querySelector('.perm-check');
    cb.checked = !cb.checked;
    item.classList.toggle('checked', cb.checked);
    chk.innerHTML = cb.checked ? '<i class="bi bi-check"></i>' : '';
}

/* Password visibility toggles */
function makeToggle(btnId, inputId) {
    const btn = document.getElementById(btnId);
    const inp = document.getElementById(inputId);
    btn.addEventListener('click', () => {
        const hidden = inp.type === 'password';
        inp.type = hidden ? 'text' : 'password';
        btn.querySelector('i').className = hidden ? 'bi bi-eye-slash' : 'bi bi-eye';
    });
}
makeToggle('btn-toggle-reader', 'reader-password');
makeToggle('btn-toggle-owner', 'owner-password');

/* Random password generator */
function randomPassword(len = 18) {
    const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789!@#$%&*';
    let out = '';
    const arr = new Uint8Array(len);
    crypto.getRandomValues(arr);
    arr.forEach(v => out += chars[v % chars.length]);
    return out;
}

document.getElementById('btn-regen-reader').addEventListener('click', () => {
    const inp = document.getElementById('reader-password');
    inp.value = randomPassword();
    inp.type = 'text';
    document.getElementById('btn-toggle-reader').querySelector('i').className = 'bi bi-eye-slash';
});

document.getElementById('btn-regen-owner').addEventListener('click', () => {
    const inp = document.getElementById('owner-password');
    inp.value = randomPassword();
    inp.type = 'text';
    document.getElementById('btn-toggle-owner').querySelector('i').className = 'bi bi-eye-slash';
});

/* Owner-only checkbox */
const chkOwnerOnly = document.getElementById('chk-owner-only');
const ownerOnlyCheck = document.getElementById('owner-only-check');
const readerInput = document.getElementById('reader-password');
const readerRegen = document.getElementById('btn-regen-reader');
const readerToggle = document.getElementById('btn-toggle-reader');

document.getElementById('lbl-owner-only').addEventListener('click', () => {
    chkOwnerOnly.checked = !chkOwnerOnly.checked;
    const on = chkOwnerOnly.checked;
    ownerOnlyCheck.innerHTML = on ? '<i class="bi bi-check"></i>' : '';
    ownerOnlyCheck.style.background = on ? 'var(--accent)' : '';
    ownerOnlyCheck.style.borderColor = on ? 'var(--accent)' : '';
    readerInput.disabled = on;
    readerRegen.disabled = on;
    readerToggle.disabled = on;
    readerInput.style.opacity = on ? '0.4' : '1';
    readerRegen.style.opacity = on ? '0.4' : '1';
    readerToggle.style.opacity = on ? '0.4' : '1';
});


document.getElementById('btn-open-file').addEventListener('click', async () => {
    const inputPath = document.getElementById('input-file-path');
    try {
        const path = await OpenFileDialog();
        inputPath.value = path;
        console.log('Input Path:', path);
    } catch (err) {
        console.error(err);
    }
});

document.getElementById('btn-save-file').addEventListener('click', async () => {
    const outputPath = document.getElementById('output-file-path');
    try {
        const path = await SaveFileDialog();
        outputPath.value = path;
        console.log('Output Path:', path);
    } catch (err) {
        console.error(err);
    }
});

/* Protect button */
document.getElementById('btn-protect').addEventListener('click', () => {
    const inputPath = document.getElementById('input-file-path').value.trim();
    const outputPath = document.getElementById('output-file-path').value.trim();
    const readerPwd = document.getElementById('reader-password').value;
    const ownerPwd = document.getElementById('owner-password').value;
    const ownerOnly = document.getElementById('chk-owner-only').checked;
    const perms = [...document.querySelectorAll('#permissions-grid input[type=checkbox]:checked')]
        .map(cb => cb.value);

    const payload = { inputPath, outputPath, readerPwd, ownerPwd, ownerOnly, perms };
    console.log('protect', payload);

    // Wails: window.go.main.App.EncryptPDF(payload).then(...).catch(...)
});


import {
    OpenFileDialog,
    SaveFileDialog,
    RandomPassword,
    EncryptPDF
} from '../wailsjs/go/gui/App';

import { EventsOn } from "../wailsjs/runtime/runtime";


/* ==========================================================================
   CONFIG
*/

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


const DOM = {
    permissionsGrid: document.getElementById('permissions-grid'),

    inputFilePath: document.getElementById('input-file-path'),
    outputFilePath: document.getElementById('output-file-path'),

    readerPassword: document.getElementById('reader-password'),
    ownerPassword: document.getElementById('owner-password'),

    btnToggleReader: document.getElementById('btn-toggle-reader'),
    btnToggleOwner: document.getElementById('btn-toggle-owner'),

    btnRegenReader: document.getElementById('btn-regen-reader'),
    btnRegenOwner: document.getElementById('btn-regen-owner'),

    btnOpenFile: document.getElementById('btn-open-file'),
    btnSaveFile: document.getElementById('btn-save-file'),

    btnProtect: document.getElementById('btn-protect'),

    chkOwnerOnly: document.getElementById('chk-owner-only'),
    ownerOnlyCheck: document.getElementById('owner-only-check'),
    lblOwnerOnly: document.getElementById('lbl-owner-only'),
};



init();

function init() {
    buildPermissionsGrid();
    bindEvents();
}


function bindEvents() {

    /* Password visibility */
    bindPasswordToggle(DOM.btnToggleReader, DOM.readerPassword);
    bindPasswordToggle(DOM.btnToggleOwner, DOM.ownerPassword);

    /* Random password */
    DOM.btnRegenReader.addEventListener('click', handleGenerateReaderPassword);
    DOM.btnRegenOwner.addEventListener('click', handleGenerateOwnerPassword);

    /* Owner only */
    DOM.lblOwnerOnly.addEventListener('click', handleOwnerOnlyToggle);

    /* File dialogs */
    DOM.btnOpenFile.addEventListener('click', handleOpenFile);
    DOM.btnSaveFile.addEventListener('click', handleSaveFile);

    /* Protect PDF */
    DOM.btnProtect.addEventListener('click', handleProtectPDF);

    /* Wails events */
    EventsOn("PDFSEC:input-selected", handleInputSelected);

    EventsOn("PDFSEC:encryption-finished", handleEncryptionFinished);
}

function handleEncryptionFinished(data){
    console.log(data);

    if (data.success) {
        console.log("PDF encrypted");
        showResultModal("PDF Encrypted!",true);
        return;
    }

    console.error(data.error);
    showResultModal(data.error,false);
}

async function handleGenerateReaderPassword() {
    const password = await generateRandomPassword();

    DOM.readerPassword.value = password;
    DOM.readerPassword.type = 'text';

    DOM.btnToggleReader.querySelector('i').className = 'bi bi-eye-slash';
}

async function handleGenerateOwnerPassword() {
    const password = await generateRandomPassword();

    DOM.ownerPassword.value = password;
    DOM.ownerPassword.type = 'text';

    DOM.btnToggleOwner.querySelector('i').className = 'bi bi-eye-slash';
}

function handleOwnerOnlyToggle() {
    DOM.chkOwnerOnly.checked = !DOM.chkOwnerOnly.checked;

    const enabled = DOM.chkOwnerOnly.checked;

    DOM.ownerOnlyCheck.innerHTML = enabled
        ? '<i class="bi bi-check"></i>'
        : '';

    DOM.ownerOnlyCheck.style.background = enabled
        ? 'var(--accent)'
        : '';

    DOM.ownerOnlyCheck.style.borderColor = enabled
        ? 'var(--accent)'
        : '';

    toggleReaderControls(enabled);
}

async function handleOpenFile() {
    try {
        const path = await OpenFileDialog();

        DOM.inputFilePath.value = path;

        console.log('Input Path:', path);

    } catch (err) {
        console.error(err);
    }
}

async function handleSaveFile() {
    try {
        const path = await SaveFileDialog();
        DOM.outputFilePath.value = path;
        console.log('Output Path:', path);

    } catch (err) {
        console.error(err);
    }
}

function handleInputSelected(data) {
    console.log("Event:", data);
    DOM.outputFilePath.value = data.outputPath;
}

async function handleProtectPDF() {
    const payload = {
        inputPath: DOM.inputFilePath.value.trim(),
        outputPath: DOM.outputFilePath.value.trim(),
        readerPwd: DOM.readerPassword.value,
        ownerPwd: DOM.ownerPassword.value,
        ownerOnly: DOM.chkOwnerOnly.checked,
        perms: getSelectedPermissions()
    };

    console.log('protect', payload);
    await EncryptPDF(payload);

}


function buildPermissionsGrid() {
    PERMISSIONS.forEach(permission => {
        const item = document.createElement('label');
        item.id = permission.id;

        item.className =
            'perm-item' +
            (permission.defaultOn ? ' checked' : '');

        item.innerHTML = `
            <input
                type="checkbox"
                name="permissions"
                value="${permission.flag}"
                ${permission.defaultOn ? 'checked' : ''}
            />

            <span class="perm-check">
                ${permission.defaultOn ? '<i class="bi bi-check"></i>' : ''}
            </span>

            <span class="perm-name">
                ${permission.label}
            </span>
        `;

        item.addEventListener('click', () => togglePermission(item));

        DOM.permissionsGrid.appendChild(item);
    });
}

function togglePermission(item) {
    const checkbox = item.querySelector('input[type=checkbox]');
    const checkIcon = item.querySelector('.perm-check');

    checkbox.checked = !checkbox.checked;

    item.classList.toggle('checked', checkbox.checked);

    checkIcon.innerHTML = checkbox.checked
        ? '<i class="bi bi-check"></i>'
        : '';
}

function bindPasswordToggle(button, input) {
    button.addEventListener('click', () => {

        const hidden = input.type === 'password';

        input.type = hidden ? 'text' : 'password';

        button.querySelector('i').className =
            hidden
                ? 'bi bi-eye-slash'
                : 'bi bi-eye';
    });
}

function toggleReaderControls(disabled) {
    DOM.readerPassword.disabled = disabled;
    DOM.btnRegenReader.disabled = disabled;
    DOM.btnToggleReader.disabled = disabled;

    DOM.readerPassword.style.opacity = disabled ? '0.4' : '1';
    DOM.btnRegenReader.style.opacity = disabled ? '0.4' : '1';
    DOM.btnToggleReader.style.opacity = disabled ? '0.4' : '1';
}


async function generateRandomPassword() {
    try {
        const password = await RandomPassword();
        return password;
    } catch (err) {
        console.error(err);
        return "";
    }
}

function getSelectedPermissions() {
    return [
        ...document.querySelectorAll(
            '#permissions-grid input[type=checkbox]:checked'
        )
    ].map(cb => cb.value);
}

function showResultModal(message, success = true) {
    if(success) alert(message);
    else alert("FAIL: " + message);
}
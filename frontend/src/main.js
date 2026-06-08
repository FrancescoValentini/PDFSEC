import {
    OpenFileDialog,
    SaveFileDialog,
    RandomPassword,
    EncryptPDF,
    IsReady
} from '../wailsjs/go/gui/App';

import { EventsOn } from "../wailsjs/runtime/runtime";
import * as bootstrap from 'bootstrap';

import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap-icons/font/bootstrap-icons.css';


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

let isAppReady = false;
let isEncrypting = false;

init();

function init() {
    buildPermissionsGrid();
    bindEvents();

    updateEncryptButtonState();
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



    /* Live validation */
    DOM.inputFilePath.addEventListener('input', updateEncryptButtonState);
    DOM.outputFilePath.addEventListener('input', updateEncryptButtonState);
    DOM.readerPassword.addEventListener('input', updateEncryptButtonState);
    DOM.ownerPassword.addEventListener('input', updateEncryptButtonState);
}

function buildEncryptionRequestPayload(){
    if (DOM.btnProtect.disabled) {
        return;
    }

    const payload = {
        inputPath: DOM.inputFilePath.value.trim(),
        outputPath: DOM.outputFilePath.value.trim(),
        readerPwd: DOM.readerPassword.value.trim(),
        ownerPwd: DOM.ownerPassword.value.trim(),
        ownerOnly: DOM.chkOwnerOnly.checked,
        perms: getSelectedPermissions()
    };

    if (!payload.inputPath || !payload.outputPath) {
        showResultModal(
            "Please select input and output files.",
            false
        );
        return;
    }

    if (payload.ownerOnly) {
        if (!payload.ownerPwd) {
            showResultModal(
                "Owner password is required.",
                false
            );
            return;
        }
    } else {
        if (!payload.readerPwd || !payload.ownerPwd) {
            showResultModal(
                "Reader and owner passwords are required.",
                false
            );
            return;
        }
    }

    return payload
}

async function handleProtectPDF() {
    const payload = await buildEncryptionRequestPayload();
    if(payload == undefined){
        return
    }

    try {
        isEncrypting = true;
        updateEncryptButtonState();
        console.log('protect', payload);
        await EncryptPDF(payload);
    } catch (err) {
        console.error(err);
        showResultModal(
            err,
            false
        );
    } finally {
        isEncrypting = false;
        updateEncryptButtonState();
    }
}

function handleEncryptionFinished(data) {
    console.log(data);
    isEncrypting = false;
    updateEncryptButtonState();
    if (data.success) {
        console.log("PDF encrypted");
        showResultModal(
            "PDF encrypted successfully!",
            true
        );
        return;
    }

    console.error(data.error);
    showResultModal(
        data.error,
        false
    );
}

async function handleGenerateReaderPassword() {
    const password = await generateRandomPassword();

    DOM.readerPassword.value = password;
    DOM.readerPassword.type = 'text';

    DOM.btnToggleReader.querySelector('i').className = 'bi bi-eye-slash';

    updateEncryptButtonState();
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

    updateEncryptButtonState();
}

async function handleOpenFile() {
    try {
        const path = await OpenFileDialog();

        DOM.inputFilePath.value = path;

        console.log('Input Path:', path);
        await checkAppReady();
    } catch (err) {
        console.error(err);
    }
}

async function handleSaveFile() {
    try {
        const path = await SaveFileDialog();
        DOM.outputFilePath.value = path;
        console.log('Output Path:', path);
        await checkAppReady();
    } catch (err) {
        console.error(err);
    }
}

async function handleInputSelected(data) {
    console.log("Event:", data);
    DOM.outputFilePath.value = data.outputPath;
    await checkAppReady();
}


function updateEncryptButtonState() {
    const inputPath = DOM.inputFilePath.value.trim();
    const outputPath = DOM.outputFilePath.value.trim();

    const isEnabled =
        isAppReady &&
        inputPath !== '' &&
        outputPath !== '';

    DOM.btnProtect.disabled = !isEnabled;

    if (isEnabled) {
        DOM.btnProtect.classList.remove('opacity-50', 'cursor-not-allowed');
    } else {
        DOM.btnProtect.classList.add('opacity-50', 'cursor-not-allowed');
    }
}

async function checkAppReady() {
    try {
        isAppReady = await IsReady();
        updateEncryptButtonState();
        return isAppReady;
    } catch (err) {
        console.error('Failed to check app ready state:', err);
        isAppReady = false;
        updateEncryptButtonState();
        return false;
    }
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

    if (disabled) DOM.readerPassword.value = '';
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
    if(success) showModal("PDF ENCRYPTED!",message,"success");
    else {
        showModal("ERROR!",message);
    }
}

/**
    Shows a Bootstrap modal with the given title and body content, centered on screen.
    Supports both error and success styles with appropriate colors.
    If a modal with the same ID already exists, it will be removed first.

    @param {string} title - The title to display in the modal header
    @param {string} body - The HTML content to display in the modal body
    @param {string} type - The type of modal: 'error' (default) or 'success'
*/
function showModal(title, body, type = 'error') {
    // Remove any existing modal with the same ID
    const existingModal = document.getElementById('dynamicModal');
    if (existingModal) existingModal.remove();

    // Define styles based on modal type
    const styles = {
        error: {
            headerClass: 'bg-danger text-white',
            borderClass: 'border-danger'
        },
        success: {
            headerClass: 'bg-success text-white',
            borderClass: 'border-success'
        }
    };

    const selectedStyle = styles[type] || styles.error;
    const modalTitle = `${title}`;

    // Create the modal HTML markup with centering classes
    const modalHTML = `
    <div class="modal fade" id="dynamicModal" tabindex="-1" aria-hidden="true">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content ${selectedStyle.borderClass}">
          <div class="modal-header ${selectedStyle.headerClass}">
            <h5 class="modal-title">${modalTitle}</h5>
          </div>
          <div class="modal-body">
            ${body}
          </div>
          <div class="modal-footer">
            <button type="button" class="btn ${type === 'error' ? 'btn-danger' : 'btn-success'}" data-bs-dismiss="modal">Close</button>
          </div>
        </div>
      </div>
    </div>`;

    // Insert the modal HTML in the page
    document.body.insertAdjacentHTML('beforeend', modalHTML);

    // Initialize and show the Bootstrap modal
    const modalElement = document.getElementById('dynamicModal');
    const modal = new bootstrap.Modal(modalElement);
    modal.show();

    // Optional: Auto-hide success modal after 3 seconds
    if (type === 'success') {
        setTimeout(() => {
            const modalInstance = bootstrap.Modal.getInstance(modalElement);
            if (modalInstance) modalInstance.hide();
        }, 3000);
    }
}

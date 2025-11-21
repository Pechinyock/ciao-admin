
const collapsedStyleName = 'collapsed'
imagePool_selected = null

function selectImagePoolItem(element) {
    if (imagePool_selected) {
        imagePool_selected.style.border = '2px solid transparent'
    }
    element.style.border = "2px solid rgb(3, 252, 78)"
    imagePool_selected = element
    const imagePoolSelectedInput = document.getElementById('selected-character-image')
    if (!imagePoolSelectedInput) {
        console.error('failed to get selected-character-image')
        return
    }
    imagePoolSelectedInput.value = element.dataset.fileName
}

function toogleCollapseButtonText(button) {
    if (button.tagName !== 'BUTTON'
        || !button.hasAttribute('data-bs-toggle')
        || button.getAttribute('data-bs-toggle') !== 'collapse') { /* RETURN */ return; }

    if (button.classList.contains(collapsedStyleName)) {
        button.textContent = 'Показать';
        button.classList.remove(collapsedStyleName);
    } else {
        button.textContent = 'Спрятать';
        button.classList.add(collapsedStyleName);
    }
}

function setFormCardType(newValue) {
    document.getElementById('creating-card-type').value = newValue
}
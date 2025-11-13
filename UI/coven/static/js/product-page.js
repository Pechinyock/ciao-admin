function productEditBtns_setOnClickEvent() {
    const productEditButtons = document.querySelectorAll('button#product-edit')
    if (!productEditButtons) {
        return
    }

    productEditButtons.forEach(btn => {
        btn.addEventListener('click', function () {
            const card = this.closest('div#product-card')
            if (!card) {
                console.error('failed to find product-card')
                return
            }
            const productId = card.dataset.productId;
            const editForm = document.getElementById('product-edit-modal-form')
            if (!editForm) {
                console.error('failed to get product edit modal')
                return
            }
            const idInput = editForm.querySelector('input#target-product-id')
            if (!idInput) {
                console.error('failed to get id input')
                return
            }

            idInput.value = productId
            const inputIdsToPassToModal = [
                "product-name",
                "product-description",
                "product-composition",
                "product-price",
            ]

            inputIdsToPassToModal.forEach(inputId => {
                const srcParagraph = card.querySelector('p#' + inputId)
                if (!srcParagraph)
                    return;
                const destinationInput = editForm.querySelector('input#' + inputId)
                if (!destinationInput) {
                    console.error('failed to find input: ' + inputId);
                    return
                }
                destinationInput.value = srcParagraph.innerText
            })
        })
    })
}

function productDeleteBtn_setOnClickEvent() {
    const productDeleteButtons = document.querySelectorAll('button#product-delete')
    if (!productDeleteButtons) {
        return
    }

    productDeleteButtons.forEach(btn => {
        btn.addEventListener('click', function () {
            const card = this.closest('div#product-card')
            if (!card) {
                console.error('failed to get card')
                return
            }
            card.classList.add('d-none')
        })
    })
}

function productModal_setOnClosingEvent(modalName) {
    const modal = document.getElementById(modalName)
    if (!modal) {
        console.error('failed to get create product modal');
        return
    }
    modal.addEventListener('hide.bs.modal', updateProductView)
}

function updateProductView() {
    htmx.ajax('GET', '/app/component/products-view', { target: '#products', swap: 'innerHTML' })
}

productEditBtns_setOnClickEvent()
productDeleteBtn_setOnClickEvent()
productModal_setOnClosingEvent('create-product-modal')
productModal_setOnClosingEvent('edit-product-modal')
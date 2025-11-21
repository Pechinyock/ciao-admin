function setModalName(element) {
    if (!element) {
        console.error('failed to set modal name the caller is null')
        return
    }
    const modalWindowContainer = document.getElementById('modal')
    if (!modalWindowContainer) {
        console.error('failed to get modal')
        return
    }
    const modalBodyLoader = document.querySelector('#body-loader')
    if (!modalBodyLoader) {
        console.error('failed modal loader')
        return
    }
    const modalName = element.dataset.modalName
    if (!modalName || modalName === "") {
        console.error('failed to get modal name');
        return
    }
    modalBodyLoader.setAttribute('hx-get', `/ui/modal-body/${modalName}`)
    const hxGetAttrValue = modalBodyLoader.getAttribute('hx-get')
    htmx.process(modalBodyLoader)
    htmx.trigger('#body-loader', 'load-modal-body')
}
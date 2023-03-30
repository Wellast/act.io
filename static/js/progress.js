const progressLockUpdate = (status) => {
    document.getElementById('progress-bar').dataset.lock = String(status);
    document.getElementById('progress-bar').hidden = !status;
}
const progressIsLock = () => document.getElementById('progress-bar').dataset.lock == 'true';

const progressCreate = () => {
    if (progressIsLock()) {
        return () => {}
    }
    progressLockUpdate(true);
    return () => progressLockUpdate(false);
}

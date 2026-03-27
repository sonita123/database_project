document.addEventListener('DOMContentLoaded', function () {
    const sidebar  = document.getElementById('sidebar');
    const backdrop = document.getElementById('sidebar-backdrop');

    function toggleSidebar() {
        const isOpen = sidebar.classList.contains('translate-x-0');

        if (isOpen) {
            sidebar.classList.add('-translate-x-full');
            sidebar.classList.remove('translate-x-0');
            backdrop.classList.add('hidden');
            document.body.style.overflow = '';
        } else {
            sidebar.classList.remove('-translate-x-full');
            sidebar.classList.add('translate-x-0');
            backdrop.classList.remove('hidden');
            document.body.style.overflow = 'hidden';
        }
    }

    // Make global
    window.toggleSidebar = toggleSidebar;

    // Backdrop click
    backdrop.addEventListener('click', toggleSidebar);

    // ESC key
    document.addEventListener('keydown', function (e) {
        if (e.key === 'Escape') {
            sidebar.classList.add('-translate-x-full');
            sidebar.classList.remove('translate-x-0');
            backdrop.classList.add('hidden');
            document.body.style.overflow = '';
        }
    });

}); // ✅ VERY IMPORTANT
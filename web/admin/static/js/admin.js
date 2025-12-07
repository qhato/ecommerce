// Admin Panel JavaScript

document.addEventListener('DOMContentLoaded', function() {
    // Sidebar submenu toggle
    const submenuToggles = document.querySelectorAll('.has-submenu > a');
    submenuToggles.forEach(toggle => {
        toggle.addEventListener('click', function(e) {
            e.preventDefault();
            const parent = this.parentElement;
            const submenu = parent.querySelector('.submenu');

            if (submenu) {
                const isOpen = submenu.classList.contains('show');

                // Close all other submenus
                document.querySelectorAll('.submenu.show').forEach(s => {
                    s.classList.remove('show');
                });

                // Toggle current submenu
                if (!isOpen) {
                    submenu.classList.add('show');
                    parent.classList.add('show');
                }
            }
        });
    });

    // Initialize tooltips if Bootstrap is loaded
    if (typeof bootstrap !== 'undefined' && bootstrap.Tooltip) {
        const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
        tooltipTriggerList.map(function (tooltipTriggerEl) {
            return new bootstrap.Tooltip(tooltipTriggerEl);
        });
    }

    // Confirm delete actions
    const deleteButtons = document.querySelectorAll('[data-confirm-delete]');
    deleteButtons.forEach(button => {
        button.addEventListener('click', function(e) {
            if (!confirm('Are you sure you want to delete this item?')) {
                e.preventDefault();
                return false;
            }
        });
    });
});

document.addEventListener('DOMContentLoaded', () => {
    // Mobile Menu Toggle
    const hamburger = document.querySelector('.hamburger');
    const navLinks = document.querySelector('.nav-links');

    if (hamburger) {
        hamburger.addEventListener('click', () => {
            navLinks.classList.toggle('active');
        });
    }

    // Typing Animation
    const typeWriter = (selector, text, speed = 100) => {
        const element = document.querySelector(selector);
        if (!element) return;
        let i = 0;
        element.innerHTML = '';
        
        function type() {
            if (i < text.length) {
                element.innerHTML += text.charAt(i);
                i++;
                setTimeout(type, speed);
            } else {
                element.style.borderRight = 'none';
            }
        }
        
        element.style.borderRight = '2px solid var(--accent-green)';
        type();
    };

    typeWriter('.typing-text', 'EPHEMERAL', 150);

    // Fade-in-up Intersection Observer
    const observerOptions = {
        threshold: 0.1
    };

    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.classList.add('visible');
            }
        });
    }, observerOptions);

    document.querySelectorAll('.fade-in-up').forEach(el => observer.observe(el));

    // GitHub API Data Fetching
    const repo = 'Mr-Dark-debug/Ephemeral';
    const fetchGitHubStats = async () => {
        try {
            const [repoRes, contributorsRes] = await Promise.all([
                fetch(`https://api.github.com/repos/${repo}`),
                fetch(`https://api.github.com/repos/${repo}/contributors`)
            ]);

            const repoData = await repoRes.json();
            const contributorsData = await contributorsRes.json();

            // Update Stats
            if (document.getElementById('star-count')) {
                document.getElementById('star-count').innerText = repoData.stargazers_count || '0';
                document.getElementById('fork-count').innerText = repoData.forks_count || '0';
                document.getElementById('issue-count').innerText = repoData.open_issues_count || '0';
            }

            // Update Contributors
            const contributorsRow = document.getElementById('contributors-row');
            if (contributorsRow) {
                contributorsRow.innerHTML = contributorsData.map(c => `
                    <a href="${c.html_url}" target="_blank" title="${c.login}">
                        <img src="${c.avatar_url}" alt="${c.login}" class="contributor-avatar">
                    </a>
                `).join('');
            }
        } catch (error) {
            console.error('Error fetching GitHub data:', error);
        }
    };

    fetchGitHubStats();

    // Copy to Clipboard
    window.copyToClipboard = (button, text) => {
        navigator.clipboard.writeText(text).then(() => {
            const originalText = button.innerText;
            button.innerText = 'COPIED!';
            setTimeout(() => {
                button.innerText = originalText;
            }, 2000);
        });
    };
});

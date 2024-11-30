let currentIndex = 0;
const items = document.querySelectorAll('.carousel-item');
const totalItems = items.length;

function moveSlide(step) {
    currentIndex += step;

    // Loop back to the beginning or end of the carousel
    if (currentIndex < 0) {
        currentIndex = totalItems - 1;
    } else if (currentIndex >= totalItems) {
        currentIndex = 0;
    }

    updateCarousel();
}

function updateCarousel() {
    const newTransformValue = -currentIndex * 100;
    document.querySelector('.carousel').style.transform = `translateX(${newTransformValue}%)`;
}

// Optional: Auto slide every 3 seconds
// setInterval(() => moveSlide(1), 3000);

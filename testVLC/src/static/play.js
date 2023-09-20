const getBuffer = async () => {
  try {
    const response = await fetch('/getAudioBuffer');
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    const buffer = await response.json();
    return buffer;
  } catch (error) {
    console.error('Error fetching buffer:', error);
    // Handle the error gracefully, e.g., display an error message to the user.
    return null;
  }
};

// Your existing code for fetching and playing audio
const play = async () => {
  const audio = document.getElementById('audioPlayer');
  console.log('Fetching audio buffer...');
  const response = await fetch('/getAudioBuffer'); // Make sure this URL is correct
  console.log('Fetch completed.');
  const arrayBuffer = await response.arrayBuffer();

  // Convert the binary data to a Blob
  const blob = new Blob([arrayBuffer], { type: 'audio/mpeg' });

  // Create a URL for the Blob
  const blobUrl = URL.createObjectURL(blob);

  // Set the audio source to the generated URL
  audio.src = blobUrl;

  // Load and play the audio
  audio.load();
  audio.play();
};

// Attach an event listener to the "Play" button
const playButton = document.getElementById('playButton');
playButton.addEventListener('click', play);

// Stop function
const stop = () => {
  const audio = document.getElementById('audioPlayer');
  audio.pause();
  audio.currentTime = 0; // Reset playback position to the beginning
};

// Attach an event listener to the "Stop" button
const stopButton = document.getElementById('stopButton');
stopButton.addEventListener('click', stop);

// Backward function (skip back by 10 seconds, adjust as needed)
const backward = () => {
  const audio = document.getElementById('audioPlayer');
  audio.currentTime -= 10; // Adjust the time as needed
};

// Attach an event listener to the "Backward" button
const backwardButton = document.getElementById('backwardButton');
backwardButton.addEventListener('click', backward);

// Forward function (skip forward by 10 seconds, adjust as needed)
const forward = () => {
  const audio = document.getElementById('audioPlayer');
  audio.currentTime += 10; // Adjust the time as needed
};

// Attach an event listener to the "Forward" button
const forwardButton = document.getElementById('forwardButton');
forwardButton.addEventListener('click', forward);

  
  
"use client";
import { useState } from 'react';

export default function Home() {
  const [file, setFile] = useState<File | null>(null);
  const [text, setText] = useState<string | null>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      setFile(e.target.files[0]);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!file) return;

    const formData = new FormData();
    formData.append('image', file);

    const res = await fetch('http://localhost:8080/upload', {
      method: 'POST',
      body: formData,
    });

    const result = await res.json();
    setText(result.text);
  };

  return (
    <div>
      <h1>Upload Image</h1>
      <form onSubmit={handleSubmit}>
        <input type="file" accept="image/*" onChange={handleFileChange} />
        <button type="submit">Upload</button>
      </form>
      {text && (
        <div>
          <h2>Extracted Text</h2>
          <p>{text}</p>
        </div>
      )}
    </div>
  );
}

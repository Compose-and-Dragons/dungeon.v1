#!/usr/bin/env python3

import os
import re
from pathlib import Path
import subprocess

# Find all Marp slide files
result = subprocess.run(
    ['find', 'deck', '-name', '*.md', '-type', 'f'],
    capture_output=True,
    text=True
)

# Filter for files with 'marp: true' and sort, excluding backup folders
slides = []
for filepath in result.stdout.strip().split('\n'):
    if not filepath:
        continue
    # Skip files in backup folders
    if '/backup/' in filepath:
        continue
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
            if re.match(r'^---\s*\nmarp:\s*true', content, re.MULTILINE):
                slides.append(filepath)
    except Exception as e:
        print(f"Error reading {filepath}: {e}")

# Sort slides by path
slides.sort()

print(f"Found {len(slides)} slide files")

# Process each slide
for i, current in enumerate(slides):
    print(f"Processing: {current}")

    # Read current file
    with open(current, 'r', encoding='utf-8') as f:
        content = f.read()

    # Remove existing navigation if present
    # Look for the navigation pattern at the end (with or without slide separator)
    nav_pattern = r'\n---\s*\n\s*\[(?:← Previous|Next →)\].*$'
    content = re.sub(nav_pattern, '', content, flags=re.MULTILINE)

    # Also remove navigation without slide separator
    nav_pattern_no_separator = r'\n\s*\[(?:← Previous|Next →)\].*$'
    content = re.sub(nav_pattern_no_separator, '', content, flags=re.MULTILINE)

    # Remove trailing whitespace
    content = content.rstrip() + '\n'

    # Prepare navigation links
    nav_parts = []

    if i > 0:
        prev = slides[i-1]
        rel_prev = os.path.relpath(prev, os.path.dirname(current))
        nav_parts.append(f"[← Previous]({rel_prev})")

    if i < len(slides) - 1:
        next_slide = slides[i+1]
        rel_next = os.path.relpath(next_slide, os.path.dirname(current))
        nav_parts.append(f"[Next →]({rel_next})")

    # Add navigation if we have links (append to last slide, not new slide)
    if nav_parts:
        nav_line = " | ".join(nav_parts)
        content += f"\n{nav_line}\n"

    # Write updated content
    with open(current, 'w', encoding='utf-8') as f:
        f.write(content)

print("Navigation added to all slide files!")

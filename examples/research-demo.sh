#!/bin/bash
# Simulates an AI agent doing research
# Run: ./examples/research-demo.sh | ./gloss watch

echo "::status id=search pending Search sources"
echo "::status id=read pending Read papers"
echo "::status id=analyze pending Analyze"
echo "::status id=write pending Write summary"
sleep 0.3

echo "::status id=search running Searching"
echo "Querying arxiv for 'transformer efficiency'..."
sleep 0.5
echo "  Found 23 papers (2024-2025)"
echo "Scanning GitHub repos..."
sleep 0.4
echo "  Found 12 relevant implementations"
echo "::status id=search done 23 papers, 12 repos"
sleep 0.2

echo "::status id=read running Reading"
echo "Reading FlashAttention-3 (Dao et al., 2024)..."
sleep 0.5
echo "Reading Ring Attention (Liu et al., 2024)..."
sleep 0.5
echo "Reading PagedAttention / vLLM..."
sleep 0.4
echo "::status id=read done 3 key papers read"
sleep 0.2

echo "::status id=analyze running Analyzing"
echo ""
echo "Key findings:"
echo "  1. FlashAttention-3: 2x speedup on H100 with FP8"
echo "  2. Ring Attention: scales to 1M+ context on 8 GPUs"
echo "  3. PagedAttention: 24x throughput via memory paging"
sleep 0.8
echo "::status id=analyze done Analysis complete"
sleep 0.2

echo "::status id=write running Writing"
sleep 0.6
echo "::status id=write done Summary ready"

echo ""
echo "::ok Research complete — 3 techniques identified"

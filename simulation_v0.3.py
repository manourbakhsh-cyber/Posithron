# ============================================================
# شبیه‌سازی فوتون (پوزیترون) — نسخه ۰.۳
# تست حملات پیچیده و مکانیزم درهم‌ریختگی (Decoherence)
# ============================================================

import numpy as np
import hashlib
from dataclasses import dataclass, field
from typing import List
import matplotlib.pyplot as plt

# === پارامترهای شبکه ===
ALPHA = 0.5
BETA  = 0.5
GAMMA = 0.3
KAPPA = 0.05
DELTA = 0.02
THETA = 0.7
LAMBDA = 0.6
W = 144
MAX_DRIFT = 60     # حداکثر اختلاف مجاز بین زمان AI و زمان انرژی (ثانیه)

SOURCE_WEIGHT = {'solar': 1.0, 'wind': 1.0, 'battery': 0.9, 'flex': 0.7}

@dataclass
class Node:
    node_id: int
    role: str = 'hybrid'
    is_attacker: bool = False
    ai_quality: float = 0.0
    ai_compute: float = 0.0
    ai_timestamp: float = 0.0
    energy_amount: float = 0.0
    energy_source: str = 'solar'
    energy_timestamp: float = 0.0
    online_periods: int = 0
    total_reward: float = 0.0

def compute_shares(node, avg_ai, avg_energy):
    A = (node.ai_quality * node.ai_compute) / max(avg_ai, 1e-9)
    A = min(A, 2.0)
    B = (node.energy_amount * SOURCE_WEIGHT.get(node.energy_source, 1.0)) / max(avg_energy, 1e-9)
    B = min(B, 2.0)
    return A, B

def compute_phase(node, peak_time, window_duration=3600):
    dt = abs(node.ai_timestamp - peak_time)
    return np.pi * min(dt / window_duration, 1.0)

def entanglement_pairs(nodes, seed):
    valid_nodes = [n for n in nodes if (n.ai_quality > 0.1 or n.energy_amount > 1.0)]
    rng = np.random.default_rng(seed)
    idx = rng.permutation(len(valid_nodes))
    partner = {n.node_id: None for n in nodes}
    for i in range(0, len(idx) - 1, 2):
        a, b = valid_nodes[idx[i]], valid_nodes[idx[i + 1]]
        if a.role != b.role:
            partner[a.node_id] = b
            partner[b.node_id] = a
    return partner

def casimir_reward(node, total_activity, threshold=1000):
    if total_activity > threshold: return 0.0
    idle_ratio = 1 - total_activity / threshold
    presence = min(1.0, node.online_periods / W)
    return KAPPA * idle_ratio * presence

def compute_reward(node, pair_node, avg_ai, avg_energy, peak_time, total_activity):
    A, B = compute_shares(node, avg_ai, avg_energy)
    casimir = casimir_reward(node, total_activity)

    # --- دفاع ۱: حمله ترکیبی (AI جعلی) ---
    is_fake_ai = (node.ai_quality > 0) and (node.ai_quality < THETA)
    if is_fake_ai: A = 0.0

    # --- دفاع ۲: درهم‌ریختگی (دستکاری زمان) ---
    time_discrepancy = abs(node.ai_timestamp - node.energy_timestamp)
    if time_discrepancy > MAX_DRIFT:
        return casimir 

    if node.ai_quality < THETA and B < LAMBDA:
        return casimir

    phi = compute_phase(node, peak_time)
    base = ALPHA * A + BETA * B

    interference = 0.0
    if A > 0 and B > 0:
        interference = GAMMA * 2 * np.sqrt(A * B) * np.cos(phi)

    E_self = 0.0
    if pair_node is not None:
        A_p, B_p = compute_shares(pair_node, avg_ai, avg_energy)
        E_self = np.sqrt(A * B_p)

    size_penalty = 1 - DELTA * min(1.0, (A + B) / 2)
    reward = (base + interference + E_self) * size_penalty
    return max(reward, casimir)

# ============================================================
# سناریو ۱: حمله ترکیبی (انرژی زیاد + AI جعلی)
# ============================================================
def scenario_hybrid_attack():
    print("\n" + "="*60)
    print("سناریو ۱: حمله ترکیبی (انرژی ۵۱٪ + AI جعلی)")
    print("="*60)
    n_nodes, n_att = 1000, 100
    nodes = [Node(i, role='hybrid') for i in range(n_nodes)]
    attackers = [Node(i+10000, role='hybrid', is_attacker=True) for i in range(n_att)]
    nodes += attackers
    
    for a in attackers:
        a.energy_amount = (0.51 * 100000) / n_att
        a.ai_quality = 0.1
        a.ai_compute = 5.0
        a.ai_timestamp, a.energy_timestamp = 0, 0
        
    for n in nodes[:n_nodes]:
        n.energy_amount = (0.49 * 100000) / n_nodes
        n.ai_quality = np.random.uniform(0.7, 1.0)
        n.ai_compute = np.random.uniform(1, 5)
        n.ai_timestamp, n.energy_timestamp = np.random.uniform(-1800, 1800), np.random.uniform(-1800, 1800)
        n.energy_timestamp = n.ai_timestamp

    h_rewards, a_rewards = [], []
    for period in range(100):
        for n in nodes: n.online_periods += 1
        t_ai = sum(n.ai_quality * n.ai_compute for n in nodes)
        t_en = sum(n.energy_amount for n in nodes)
        avg_ai, avg_en = t_ai/len(nodes), t_en/len(nodes)
        pairs = entanglement_pairs(nodes, seed=period)
        h_r, a_r = 0.0, 0.0
        for n in nodes:
            r = compute_reward(n, pairs.get(n.node_id), avg_ai, avg_en, 0, t_ai+t_en)
            if n.is_attacker: a_r += r
            else: h_r += r
        h_rewards.append(h_r/n_nodes)
        a_rewards.append(a_r/n_att)

    avg_h, avg_a = np.mean(h_rewards[-20:]), np.mean(a_rewards[-20:])
    print(f"پاداش صادق: {avg_h:.4f} | پاداش مهاجم: {avg_a:.4f}")
    print(f"نسبت مهاجم/صادق: {avg_a/avg_h:.2f}x")
    print("✅ دفاع موفق" if avg_a/avg_h < 1.2 else "❌ دفاع ناموفق")
    return h_rewards, a_rewards

# ============================================================
# سناریو ۲: دستکاری اوراکل فاز (Oracle Manipulation)
# ============================================================
def scenario_oracle_manipulation():
    print("\n" + "="*60)
    print("سناریو ۲: دستکاری زمان برای هم‌زمانی جعلی")
    print("="*60)
    n_nodes, n_att = 500, 50
    nodes = [Node(i) for i in range(n_nodes)] + [Node(i+5000, is_attacker=True) for i in range(n_att)]
    
    for n in nodes:
        n.ai_quality = 0.9
        n.ai_compute = 3.0
        n.energy_amount = 10.0
        if n.is_attacker:
            n.ai_timestamp = 0 
            n.energy_timestamp = np.random.uniform(500, 1800)
        else:
            n.ai_timestamp = np.random.uniform(-1800, 1800)
            n.energy_timestamp = n.ai_timestamp

    h_rewards, a_rewards = [], []
    for period in range(100):
        for n in nodes: n.online_periods += 1
        t_ai = sum(n.ai_quality * n.ai_compute for n in nodes)
        t_en = sum(n.energy_amount for n in nodes)
        avg_ai, avg_en = t_ai/len(nodes), t_en/len(nodes)
        pairs = entanglement_pairs(nodes, seed=period)
        h_r, a_r = 0.0, 0.0
        for n in nodes:
            r = compute_reward(n, pairs.get(n.node_id), avg_ai, avg_en, 0, t_ai+t_en)
            if n.is_attacker: a_r += r
            else: h_r += r
        h_rewards.append(h_r/n_nodes)
        a_rewards.append(a_r/n_att)

    avg_h, avg_a = np.mean(h_rewards[-20:]), np.mean(a_rewards[-20:])
    print(f"پاداش صادق: {avg_h:.4f} | پاداش مهاجم: {avg_a:.4f}")
    print(f"نسبت مهاجم/صادق: {avg_a/avg_h:.2f}x")
    print("✅ دفاع موفق: درهم‌ریختگی اعمال شد" if avg_a/avg_h < 0.5 else "❌ دفاع ناموفق")
    return h_rewards, a_rewards

# ============================================================
# سناریو ۳: حمله سایبیل در درهم‌تنیدگی
# ============================================================
def scenario_sybil_entanglement():
    print("\n" + "="*60)
    print("سناریو ۳: حمله سایبیل (ساخت نودهای فیک)")
    print("="*60)
    n_nodes, n_att_main, n_sybils = 200, 20, 500
    nodes = [Node(i) for i in range(n_nodes)]
    attackers = [Node(i+1000, is_attacker=True, ai_quality=0.9, ai_compute=5.0, energy_amount=20.0) for i in range(n_att_main)]
    sybils = [Node(i+2000, is_attacker=True, ai_quality=0.0, ai_compute=0.0, energy_amount=0.0) for i in range(n_sybils)]
    nodes += attackers + sybils
    
    for n in nodes[:n_nodes]:
        n.ai_quality, n.ai_compute, n.energy_amount = 0.8, 2.0, 10.0
        n.ai_timestamp, n.energy_timestamp = 0, 0

    h_rewards, a_rewards = [], []
    for period in range(100):
        for n in nodes: n.online_periods += 1
        t_ai = sum(n.ai_quality * n.ai_compute for n in nodes)
        t_en = sum(n.energy_amount for n in nodes)
        avg_ai, avg_en = t_ai/len(nodes), t_en/len(nodes)
        pairs = entanglement_pairs(nodes, seed=period)
        h_r, a_r = 0.0, 0.0
        for n in nodes:
            r = compute_reward(n, pairs.get(n.node_id), avg_ai, avg_en, 0, t_ai+t_en)
            if n.is_attacker: a_r += r
            else: h_r += r
        h_rewards.append(h_r/n_nodes)
        a_rewards.append(a_r/n_att_main)

    avg_h, avg_a = np.mean(h_rewards[-20:]), np.mean(a_rewards[-20:])
    print(f"پاداش صادق: {avg_h:.4f} | پاداش مهاجم اصلی: {avg_a:.4f}")
    print(f"نسبت مهاجم/صادق: {avg_a/avg_h:.2f}x")
    print("✅ دفاع موفق" if avg_a/avg_h < 1.5 else "❌ دفاع ناموفق")
    return h_rewards, a_rewards

# ============================================================
# اجرا
# ============================================================
if __name__ == "__main__":
    np.random.seed(42)
    
    h1, a1 = scenario_hybrid_attack()
    h2, a2 = scenario_oracle_manipulation()
    h3, a3 = scenario_sybil_entanglement()

    fig, axes = plt.subplots(1, 3, figsize=(15, 4))
    axes[0].plot(h1, label='صادق', alpha=0.7); axes[0].plot(a1, label='مهاجم ترکیبی', alpha=0.7)
    axes[0].set_title('سناریو ۱: حمله ترکیبی'); axes[0].legend()
    axes[1].plot(h2, label='صادق', alpha=0.7); axes[1].plot(a2, label='مهاجم اوراکل', alpha=0.7)
    axes[1].set_title('سناریو ۲: دستکاری زمان'); axes[1].legend()
    axes[2].plot(h3, label='صادق', alpha=0.7); axes[2].plot(a3, label='مهاجم سایبیل', alpha=0.7)
    axes[2].set_title('سناریو ۳: حمله سایبیل'); axes[2].legend()
    plt.tight_layout()
    plt.savefig('posithron_v3_simulation.png', dpi=100)
    print("\n📊 نمودار ذخیره شد: posithron_v3_simulation.png")

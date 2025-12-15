# ☁️ Proxmox Cloud-Init Template (Ubuntu 24.04) — Full Step

สร้าง Ubuntu 24.04 Cloud-Init Template บน Proxmox เพื่อใช้ Clone VM ได้รวดเร็ว
พร้อมตั้งค่า Cloud-Init ให้ครบ:
- User / Password
- SSH public key
- IPv4 + Gateway
- (Optional) DNS
- Upgrade packages

---

## ✅ Prerequisites

- Proxmox node เข้าถึงได้
- Storage มี `local-lvm`
- มี SSH public key สำหรับ Ansible (เช่น `ansible-key.pub`)

> NOTE: ใน Proxmox Web UI ปกติ SSH เข้าที่ port `22`  
> ถ้าเข้าผ่านเว็บคือ `https://<IP>:8006` แต่ SSH ไม่ใช่ `-p 8006`

---

## 1) SSH เข้า Proxmox Node

```bash
ssh root@10.13.104.216

wget https://cloud-images.ubuntu.com/releases/noble/release/ubuntu-24.04-server-cloudimg-amd64.img \
  -O ubuntu-24.04-cloudimg-amd64.img

qm create 9000 \
  --name ubuntu-2404-cloudinit \
  --memory 16384 \
  --cores 8 \
  --net0 virtio,bridge=vmbr0

qm importdisk 9000 ubuntu-24.04-cloudimg-amd64.img local-lvm

qm set 9000 --scsihw virtio-scsi-pci --scsi0 local-lvm:vm-9000-disk-0
qm set 9000 --ide2 local-lvm:cloudinit
qm set 9000 --boot c --bootdisk scsi0
qm set 9000 --serial0 socket --vga serial0
qm set 9000 --agent enabled=1

ls -la /root/ansible-key.pub

ssh-keygen -t ed25519 -f /root/ansible-key -N ""
# จะได้:
# /root/ansible-key      (private key)
# /root/ansible-key.pub  (public key)


qm set 9000 \
  --ciuser dev \
  --cipassword 'YOUR_PASSWORD' \
  --sshkeys /root/ansible-key.pub \
  --ipconfig0 ip=10.13.104.150/24,gw=10.13.104.254 \
  --ciupgrade 1

qm template 9000
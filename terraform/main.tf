resource "proxmox_virtual_environment_vm" "k3s_nodes" {
  # define 3 nodes, VM ids, and static IPs
  for_each = {
    "k3s-master"   = { vmid = 801, ip = "10.0.0.210" }
    "k3s-worker-1" = { vmid = 802, ip = "10.0.0.211" }
    "k3s-worker-2" = { vmid = 803, ip = "10.0.0.212" }
  }

  node_name = "proxmox"
  name      = each.key
  vm_id     = each.value.vmid

  # clone from template
  clone {
    vm_id = 8000
  }

  # hardware resources (4gb ram, 2 cores per node)
  cpu {
    cores = 2
    type  = "host"
  }

  memory {
    dedicated = 4096
  }

  # expand the template disk to 30gb for k8s container images
  disk {
    datastore_id = "local-lvm"
    interface    = "scsi0"
    size         = 30
  }

  # enable the qemu guest agent so proxmox can monitor VMs
  agent {
    enabled = false
  }

  # Cloud-Init config
  initialization {
    ip_config {
      ipv4 {
        address = "${each.value.ip}/24"
        gateway = "10.0.0.1"
      }
    }
    user_account {
      username = "ubuntu"
      keys     = [var.ssh_public_key]
    }
  }
}
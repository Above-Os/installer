<template>
  <q-page class="row items-center justify-evenly">
    <div>
      <div
        v-if="
          store.install_status == 'Download' ||
          store.install_status == 'Install' ||
          store.install_status == 'Fail'
        "
      >
        <div class="column">
          <div class="row">Status: {{ store.install_status }}</div>
          <q-scroll>
            <div v-for="(message, index) in store.msg" :key="index">
              <div class="row">
                {{ convertMilliseconds(message.time) }} {{ message.info }}
              </div>
            </div>
          </q-scroll>
        </div>
      </div>
      <div v-else-if="store.install_status == 'Not_Started'">
        <q-form @submit="onInstall" class="q-gutter-md">
          <q-input
            filled
            v-model="name"
            label="Terminus Name *"
            hint="Terminus Name"
            lazy-rules
            :rules="[
              (val) => (val && val.length > 0) || 'Please type something',
            ]"
          />

          <q-btn
            class="btn-size-xs text-caption q-ml-lg text-grey-8 copy-button"
            color="grey-2"
            label="install"
            outline
            no-caps
            type="submit"
          />
        </q-form>
      </div>
      <div v-else></div>
    </div>
  </q-page>
</template>

<script lang="ts" setup>
import { useInstallStore } from 'stores/install';
//import { useQuasar } from 'quasar';
import { ref, onMounted, onUnmounted } from 'vue';
import { BtNotify, NotifyDefinedType } from '@bytetrade/ui';

const store = useInstallStore();
const name = ref('');
let timer: any;

function convertMilliseconds(ms: any) {
  // 创建一个Date对象，传入毫秒数
  const date = new Date(ms);

  // 获取小时、分钟和秒数
  const hours = date.getUTCHours();
  const minutes = date.getUTCMinutes();
  const seconds = date.getUTCSeconds();

  // 格式化输出（确保两位数显示）
  const formattedHours = String(hours).padStart(2, '0');
  const formattedMinutes = String(minutes).padStart(2, '0');
  const formattedSeconds = String(seconds).padStart(2, '0');

  return `${formattedHours}:${formattedMinutes}:${formattedSeconds}`;
}

async function onInstall() {
  console.log(name.value);
  if (name.value.indexOf('@') < 0) {
    BtNotify.show({
      type: NotifyDefinedType.FAILED,
      message: 'Please input a terminus name',
    });
    return;
  }

  await store.install({
    config: {
      terminus_os_domainname: name.value.split('@')[1], // 域名（可为空，空的话会设置默认值 myterminus.com）
      terminus_os_username: name.value.split('@')[0], // 用户名（必填项）
      kube_type: 'k8s', // 安装类型（可为空，空的话会设置默认值 k3s）；选项值：k8s、k3s
      vendor: 'private', // 所属位置（可为空，空的话会设置默认值 private）；选项值：private、aws、aliyun
      gpu_enable: 1, // 是否开启 gpu（必填项，0 或 1）；暂时未用到
      gpu_share: 1, // gpu 共享（必填项，0 或 1）；暂时未用到
      version: '',
    },
  });
}

onMounted(() => {
  timer = setInterval(async () => {
    if (
      store.install_status == 'Download' ||
      store.install_status == 'Install'
    ) {
      await store.status();
    }
  }, 1000);
});

onUnmounted(() => {
  if (timer) {
    clearInterval(timer);
  }
});
</script>

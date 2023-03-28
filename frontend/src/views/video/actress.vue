<template>
  <div>
    <a-card>
      <a-form-item label="管理模式">
        <a-switch v-model:checked="manageMode"/>
      </a-form-item>
      <a-form :model="searchValue" layout="inline" autocomplete="off">
        <a-form-item label="演员名" name="name">
          <a-input v-model:value="searchValue.name"/>
        </a-form-item>

        <a-form-item>
          <a-button type="primary" @click="() => {refresh()}" >搜索</a-button>
        </a-form-item>
      </a-form>
      <a-divider/>
      <a-button-group>
        <a-button type="primary" @click="create">新建演员</a-button>
        <a-popconfirm v-if="manageMode" title="确认删除" @confirm="deleteSelect">
          <a-button type="danger">删除选中</a-button>
        </a-popconfirm>
      </a-button-group>

      <a-divider/>
      <a-table v-if="manageMode" :dataSource="actresses"
               :row-selection="{ selectedRowKeys: selectedRowKeys, onChange: onSelectChange }" rowKey="id" :columns="columns">
        <template #bodyCell="{ column, text, record }">
          <template v-if="column.dataIndex === 'star'">
            <a-switch disabled v-model:checked="record.star"/>
          </template>
          <template v-else-if="column.dataIndex === 'action'">
            <a-button size="middle" type="primary" @click="updateRecord(record)">
              <template #icon>
                <EditOutlined/>
              </template>
            </a-button>
            <a-popconfirm title="确认删除" @confirm="deleteRecord(record.id)">
              <a-button size="middle" type="danger">
                <template #icon>
                  <DeleteOutlined/>
                </template>
              </a-button>
            </a-popconfirm>
          </template>
        </template>
      </a-table>
      <a-card v-else>
        <template v-for="(actress, index) in actresses" :key="index">
          <a-tag :color="randomTagColor()" style="margin-top: 5px" @click="() => {

          }">
            {{ actress.name }}
          </a-tag>
        </template>
      </a-card>
    </a-card>
    <a-modal v-model:visible="visible" :title="modalTitle" @ok="submit">
      <a-form :model="actress" :label-col="{ style: { width: '150px' } }" :wrapper-col="{ span: 14 }">
        <a-form-item label="演员名" :rules="[{ required: true }]">
          <a-input v-model:value="actress.name" placeholder="演员名 “,” 分隔可创建多个"/>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, ref} from 'vue';
import {createActress, deleteActress, listActress, updateActress} from "@/api/actress";
import {DeleteOutlined, EditOutlined} from '@ant-design/icons-vue';
import {message} from 'ant-design-vue';
import {randomTagColor} from "@/utils/color";

const manageMode = ref<boolean>(false);
const visible = ref<boolean>(false);
const modalTitle = ref<string>('');

const empty = {
  id: 0,
  name: '',
}

const selectedRowKeys = ref<any[]>([]);
const onSelectChange = (rowKeys: any[]) => {
  selectedRowKeys.value = rowKeys;
};

const actress = ref<any>(empty);

const create = () => {
  visible.value = true;
  modalTitle.value = '添加演员';
  actress.value = {...empty};
};

const submit = () => {
  if (!actress.value.name) {
    message.warning('请输入演员名');
    return;
  }
  if (actress.value.id === 0) {
    createActress({
      name: actress.value.name,
    }).then(() => {
      refresh();
    });
  } else {
    updateActress({
      id: actress.value.id,
      name: actress.value.name,
    }).then(() => {
      refresh();
    })
  }
  refresh();
  visible.value = false;
};

const updateRecord = (record: any) => {
  visible.value = true;
  modalTitle.value = '修改演员';
  actress.value = {
    id: record.id,
    name: record.name,
  };
};

const deleteRecord = (id: number) => {
  deleteActress({
    id: [id]
  }).then(() => {
    refresh();
  });
};
const deleteSelect = () => {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择要删除的演员');
    return;
  }
  deleteActress({
    id: selectedRowKeys.value
  }).then(() => {
    refresh();
  });
};


const columns = [
  {
    title: 'id',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '演员名',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '创建时间',
    dataIndex: 'create_time',
    key: 'create_time',
  },
  {
    title: '操作',
    dataIndex: 'action',
  },
];
const actresses = ref<any[]>([]);

const searchValue = ref<any>({
  name: '',
});
const refresh = () => {
  const data = {
    name: searchValue.value.name,
  };
  listActress(data).then((res) => {
    actresses.value = res.data.data;
  });
};

onMounted(() => {
  refresh();
})

</script>
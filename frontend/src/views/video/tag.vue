<template>
  <a-card>
    <a-form-item label="管理模式">
        <a-switch v-model:checked="manageMode"/>
    </a-form-item>
    <a-form :model="searchValue" layout="inline" autocomplete="off">
      <a-form-item label="标签" name="name">
        <a-input v-model:value="searchValue.name"/>
      </a-form-item>

      <a-form-item>
        <a-button type="primary" @click="() => {refresh()}" >搜索</a-button>
      </a-form-item>
    </a-form>
    <a-divider/>
    <a-button-group>
      <a-popover v-model:visible="visible" title="添加标签" trigger="click">
        <template #content>
          <a-input v-model:value="tagValue" placeholder="标签名 “,” 分隔可创建多个"/>
          <a-button style="margin-top: 10px" size="middle" @click="ensure">确定</a-button>
        </template>
      <a-button type="primary">添加标签</a-button>
      </a-popover>
      <a-popconfirm v-if="manageMode" title="确认删除" @confirm="deleteSelect">
        <a-button type="danger">删除选中</a-button>
      </a-popconfirm>
    </a-button-group>

    <a-divider/>
    <a-table v-if="manageMode" :dataSource="tags"
             :row-selection="{ selectedRowKeys: selectedRowKeys, onChange: onSelectChange }" rowKey="id" :columns="columns">
      <template #bodyCell="{ column, text, record }">
        <template v-if="column.dataIndex === 'action'">
          <a-popconfirm title="确认删除" @confirm="deleteRecord(record.id)">
            <a-button size="middle" danger>
              <template #icon>
                <DeleteOutlined/>
              </template>
            </a-button>
          </a-popconfirm>
        </template>
      </template>
    </a-table>
    <a-card v-else>
      <template v-for="(tag, index) in tags" :key="index">
        <a-tag :color="randomTagColor()" style="margin-top: 5px" @click="() => {

          }">
          {{ tag.name }}
        </a-tag>
      </template>
    </a-card>
  </a-card>
</template>

<script lang="ts" setup>
import {onMounted, ref} from 'vue';
import {createTag, deleteTag, listTag} from "@/api/tag";
import {message} from 'ant-design-vue';
import {DeleteOutlined} from '@ant-design/icons-vue';
import {randomTagColor} from "@/utils/color";

const manageMode = ref<boolean>(false);

const tagValue = ref<string>('');
const visible = ref<boolean>(false);
const ensure = () => {
  if (!tagValue.value) {
    message.warning('请输入标签名');
    return;
  }
  createTag({
    name: tagValue.value
  }).then(() => {
    visible.value = false;
    tagValue.value = '';
    refresh();
  });
  refresh();
  visible.value = false;
};

const deleteRecord = (id: number) => {
  deleteTag({
    id: [id]
  }).then(() => {
    refresh();
  });
};
const deleteSelect = () => {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择要删除的标签');
    return;
  }
  deleteTag({
    id: selectedRowKeys.value
  }).then(() => {
    refresh();
  });
};

const selectedRowKeys = ref<any[]>([]);
const onSelectChange = (rowKeys: any[]) => {
  selectedRowKeys.value = rowKeys;
};


const columns = [
  {
    title: 'id',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '标签名',
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
const tags = ref<any[]>([]);
const searchValue = ref<any>({
  name: '',
});

const refresh = () => {
  const data = {
    name: searchValue.value.name,
  };
  listTag(data).then((res) => {
    tags.value = res.data.data;
  });
};

onMounted(() => {
  refresh();
})

</script>
<template>
  <div>
    <a-card>
      <a-button type="primary" @click="create">新建根节点</a-button>
      <a-divider/>
      <a-tree
          block-node
          v-model:expandedKeys="allKeys"
          :tree-data="tree"
          :field-names="{
            title: 'name',
            key: 'id',
            children: 'children',
          }">
        <template #title="{ id, name, children }">
          <span>
            {{ name }}
            <a-button-group style="float: right; margin-right: 50px">
              <a-button size="small" @click="addChildNode(id)">
                <template #icon>
                  <PlusOutlined/>
                </template>
              </a-button>
              <a-button size="small" @click="editChildNode(id, name)">
                <template #icon>
                  <EditOutlined/>
                </template>
              </a-button>
              <a-popconfirm v-if="children?.length === 0" title="确认删除" @confirm="deleteNode(id)">
                <a-button  size="small" danger>
                  <template #icon>
                    <DeleteOutlined/>
                  </template>
                </a-button>
              </a-popconfirm>
            </a-button-group>
          </span>
        </template>
      </a-tree>
    </a-card>
    <a-modal v-model:visible="visible" :title="modalTitle" @ok="submit">
      <a-form :model="node" :label-col="{ style: { width: '150px' } }" :wrapper-col="{ span: 14 }">
        <a-form-item label="名称" :rules="[{ required: true }]">
          <a-input v-model:value="node.name"/>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>
<script lang="ts" setup>
import {onMounted, ref} from 'vue';
import {message} from 'ant-design-vue';
import {createTreeNode, deleteTreeNode, listTree, updateTreeNode} from "@/api/tree";
import {DeleteOutlined, EditOutlined, PlusOutlined} from '@ant-design/icons-vue';

const visible = ref<boolean>(false);
const modalTitle = ref<string>('');

const empty = {
  id: 0,
  name: '',
  parent_id: 0,
};

const node = ref<any>(empty);

const create = () => {
  visible.value = true;
  modalTitle.value = '添加根节点';
  node.value = {...empty};
};

const submit = () => {
  if (!node.value.name) {
    message.warning('请输入名称');
    return;
  }
  if (node.value.id === 0) {
    createTreeNode({
      name: node.value.name,
      parent_id: node.value.parent_id,
    }).then(() => {
      refresh();
    });
  } else {
    updateTreeNode({
      id: node.value.id,
      name: node.value.name,
      parent_id: node.value.parent_id,
    }).then(() => {
      refresh();
    });
  }
  visible.value = false;
};

const addChildNode = (id: number) => {
  visible.value = true;
  modalTitle.value = '添加子节点';
  node.value = {...empty, parent_id: id};
};

const editChildNode = (id: number, name: string) => {
  visible.value = true;
  modalTitle.value = '编辑节点';
  node.value = {
    id: id,
    name: name,
  };
};

const deleteNode = (id: number) => {
  deleteTreeNode({
    id: id
  }).then(() => {
    refresh();
  });
};


const tree = ref<any[]>([]);

const refresh = () => {
  const data = {};
  listTree(data).then((res) => {
    tree.value = buildTree(res.data.data);
  });
};

const allKeys = ref<number[]>([]);

const buildTree = (data: any[]) => {
  allKeys.value = [];
  const tree: any[] = []
  const treeMap = {};
  data.forEach((item) => {
    allKeys.value.push(item.id);
    treeMap[item.id] = item;
    item.children = [];
    if (item.parent_id === 0) {
      tree.push(item)
    }
  })
  data.forEach((item) => {
    if (item.parentId !== 0) {
      treeMap[item.parent_id]?.children.push(item)
    }
  })
  return tree
}

onMounted(() => {
  refresh();
})


</script>

<style>
.custom-tree-node {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 16px;
  padding-right: 8px;
}

</style>



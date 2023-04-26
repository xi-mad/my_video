<template>
  <div>
    <a-card>
      <a-form-item label="管理模式">
        <a-switch v-model:checked="manageMode"/>
      </a-form-item>
      <a-form :model="searchValue" layout="inline" autocomplete="off">

        <a-form-item label="合集名称" name="name">
          <a-input v-model:value="searchValue.name"/>
        </a-form-item>
        <a-form-item label="演员" name="actress">
          <a-select
              style="width: 200px"
              show-search
              mode="multiple"
              v-model:value="searchValue.actress"
              :filter-option="filterOption"
              :options="actressOptions"
          ></a-select>
        </a-form-item>
        <a-form-item label="标签" name="tag">
          <a-select
              style="width: 200px"
              v-model:value="searchValue.tag"
              mode="multiple"
              :filter-option="filterOption"
              :options="tagOptions"
          ></a-select>
        </a-form-item>
        <a-form-item label="分类" name="tree">
          <a-tree-select
              style="width: 200px"
              v-model:value="searchValue.tree"
              show-search
              :dropdown-style="{ maxHeight: '400px', overflow: 'auto' }"
              allow-clear
              multiple
              :show-checked-strategy="SHOW_ALL"
              tree-default-expand-all
              :tree-data="treeOptions"
          >
            <template #tagRender="{ label, closable, onClose, option }">
              <a-tag :closable="closable" :color="option.color" style="margin-right: 3px" @close="onClose">
                {{ label }}&nbsp;&nbsp;
              </a-tag>
            </template>
          </a-tree-select>

        </a-form-item>

        <a-form-item>
          <a-button type="primary" @click="() => {refresh()}" >搜索</a-button>
        </a-form-item>
      </a-form>
      <a-divider/>
      <a-button-group>
        <a-button type="primary" @click="create">新建合集</a-button>
        <a-popconfirm v-if="manageMode" title="确认删除" @confirm="deleteSelect">
          <a-button type="danger">删除选中</a-button>
        </a-popconfirm>
      </a-button-group>
      <a-divider/>
      <a-pagination
          v-model:current="page"
          v-model:pageSize="pageSize"
          v-model:total="total"
          show-size-changer
          @change="sizeChange"
      />
      <div v-if="manageMode">
        <a-image-preview-group>
          <a-table :dataSource="objects"
                   :row-selection="{ selectedRowKeys: selectedRowKeys, onChange: onSelectChange }" rowKey="id" :columns="columns" :pagination="false">
            <template #bodyCell="{ column, text, record }">
              <template v-if="column.dataIndex === 'thumbnail'">
                <a-image :src="'data:image/jpg;base64,' + record.thumb"
                         :alt="record.name"
                         :width="100"
                />
              </template>
              <template v-else-if="column.dataIndex === 'description'">
                <a href="#" :title="record.description">{{record.description.substring(0, 20)}}</a>
              </template>
              <template v-else-if="column.dataIndex === 'action'">
                <a-button size="small" type="primary" @click="updateRecord(record)">
                  <template #icon>
                    <EditOutlined/>
                  </template>
                </a-button>
                <a-popconfirm title="确认删除" @confirm="deleteRecord(record.id)">
                  <a-button size="small" type="danger">
                    <template #icon>
                      <DeleteOutlined/>
                    </template>
                  </a-button>
                </a-popconfirm>
              </template>
            </template>
          </a-table>
        </a-image-preview-group>
      </div>
      <div v-else>
        <a-card>
          <a-image-preview-group>
            <ul class="custom_container">
              <template v-for="(object, index) in objects" :key="index">
                <li class="custom_card">
                  <a-image :src="'data:image/jpg;base64,' + object.thumb"/>
                  <a-row style="margin-top: 5px"><p style="text-overflow: ellipsis; overflow: hidden; margin-bottom: 0">{{ object.name }}</p></a-row>
                </li>
              </template>
            </ul>
          </a-image-preview-group>
        </a-card>
      </div>
      <a-pagination
          v-model:current="page"
          v-model:pageSize="pageSize"
          v-model:total="total"
          show-size-changer
          @change="sizeChange"
      />
    </a-card>
    <a-modal v-model:visible="visible" :title="modalTitle" @ok="submit" :maskClosable="false" :keyboard="false">
      <a-form :model="object" :label-col="{ style: { width: '150px' } }" :wrapper-col="{ span: 14 }">
        <a-form-item label="路径" :rules="[{ required: true }]">
          <a-input v-model:value="object.path"/>
        </a-form-item>
        <a-form-item label="名称">
          <a-input v-model:value="object.name"/>
        </a-form-item>
        <a-form-item label="演员">
          <a-select
              v-model:value="object.actress"
              mode="multiple"
              style="width: 100%"
              :filter-option="filterOption"
              :options="actressOptions"
          ></a-select>
        </a-form-item>
        <a-form-item label="标签">
          <a-select
              v-model:value="object.tag"
              mode="multiple"
              style="width: 100%"
              :filter-option="filterOption"
              :options="tagOptions"
          ></a-select>
        </a-form-item>
        <a-form-item label="分组">
          <a-tree-select
              v-model:value="object.tree"
              show-search
              style="width: 100%"
              :dropdown-style="{ maxHeight: '400px', overflow: 'auto' }"
              allow-clear
              multiple
              :show-checked-strategy="SHOW_ALL"
              tree-default-expand-all
              :tree-data="treeOptions"
          >
            <template #tagRender="{ label, closable, onClose, option }">
              <a-tag :closable="closable" :color="option.color" style="margin-right: 3px" @close="onClose">
                {{ label }}&nbsp;&nbsp;
              </a-tag>
            </template>
          </a-tree-select>
        </a-form-item>
        <a-form-item label="磁力链接">
          <a-input v-model:value="object.magnet"/>
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea
              v-model:value="object.description"
              placeholder=""
              :auto-size="{ minRows: 2, maxRows: 5 }"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, ref} from 'vue';
import {createObject, deleteObject, listObject, updateObject} from "@/api/object";
import {DeleteOutlined, EditOutlined} from '@ant-design/icons-vue';
import {message, TreeSelect} from 'ant-design-vue';
import {optionsActress} from "@/api/actress";
import {optionsTag} from "@/api/tag";
import {optionsTree} from "@/api/tree";
import {buildTree} from "@/utils/util";

const SHOW_ALL = TreeSelect.SHOW_ALL;

const manageMode = ref<boolean>(false);

const visible = ref<boolean>(false);
const modalTitle = ref<string>('');

const empty = {
  id: 0,
  name: '',
  description: '',
  path: '',
  magnet: '',
  scanPath: '',

  actress: [],
  tag: [],
  tree: [],
}

const object = ref<any>(empty);

const create = () => {
  visible.value = true;
  modalTitle.value = '添加对象';
  object.value = {...empty};
};

const submit = () => {
  let valid = true;
  if (!object.value.path) {
    message.warning('请输入路径');
    valid = false;
  }

  if (valid) {
    if (object.value.id === 0) {
      createObject({
        name: object.value.name,
        description: object.value.description,
        path: object.value.path,
        magnet: object.value.magnet,
        actress: object.value.actress,
        tag: object.value.tag,
        tree: object.value.tree,

      }).then(() => {
        refresh();
      });
    } else {
      updateObject({
        id: object.value.id,
        name: object.value.name,
        description: object.value.description,
        path: object.value.path,
        magnet: object.value.magnet,
        actress: object.value.actress,
        tag: object.value.tag,
        tree: object.value.tree,

      }).then(() => {
        refresh();
      })
    }
    visible.value = false;
  }
};

const updateRecord = (record: any) => {
  visible.value = true;
  modalTitle.value = '修改对象';
  object.value = {
    id: record.id,
    name: record.name,
    description: record.description,
    path: record.path,
    magnet: record.magnet,
    actress: record.actress,
    tag: record.tag,
    tree: record.tree,
  };
};

const deleteRecord = (id: number) => {
  deleteObject({
    id: [id]
  }).then(() => {
    refresh();
  });
};
const deleteSelect = () => {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择要删除的对象');
    return;
  }
  deleteObject({
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
    title: '编号',
    dataIndex: 'num',
    key: 'num',
  },
  {
    title: '发布日',
    dataIndex: 'release',
    key: 'release',
  },
  {
    title: '标签',
    dataIndex: 'label',
    key: 'label',
  },
  {
    title: '名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '描述',
    dataIndex: 'description',
    key: 'description',
  },
  {
    title: '路径',
    dataIndex: 'path',
    key: 'path',
  },
  {
    title: '磁力链接',
    dataIndex: 'magnet',
    key: 'magnet',
  },
  {
    title: '预览图',
    dataIndex: 'thumbnail',
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

const objects = ref<any[]>([]);
const actressOptions = ref<any[]>([]);
const tagOptions = ref<any[]>([]);
const treeOptions = ref<any[]>([]);

const treeMap = {};
const actressMap = {};
const tagMap = {};

const filterOption = (input: string, option: any) => {
  return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};

const page = ref<number>(1);
const pageSize = ref<number>(20);
const total = ref<number>(0);

const sizeChange = (nPage: number, nPageSize: number) => {
  page.value = nPage;
  pageSize.value = nPageSize;
  refresh();
};

const searchValue = ref<any>({
  name: '',
  actress: [],
  tag: [],
  tree: [],
});

const refresh = () => {
  const data = {
    page: page.value,
    page_size: pageSize.value,
    path: searchValue.value.path,
    actress: searchValue.value.actress,
    tag: searchValue.value.tag,
    tree: searchValue.value.tree,
    nfo: searchValue.value.nfo,
  };
  listObject(data).then((res) => {
    objects.value = res.data.data.data;
    total.value = res.data.data.total;
  });
};


onMounted(() => {
  optionsActress().then((res) => {
    actressOptions.value = res.data.data;
    res.data.data.forEach((item: any) => {
      actressMap[item.value] = item.label;
    });
  });
  optionsTag().then((res) => {
    tagOptions.value = res.data.data;
    res.data.data.forEach((item: any) => {
      tagMap[item.value] = item.label;
    });
  });
  optionsTree().then((res) => {
    treeOptions.value = buildTree(res.data.data);
    res.data.data.forEach((item: any) => {
      treeMap[item.value] = item.label;
    });
  });

  setTimeout(() => {
    refresh();
  }, 100);
})

</script>
<style scoped>
.custom_container {
  list-style: none;
  column-gap: 0;
  padding: 0;
  column-count: 5;
}
.custom_card {
  width: 100%;
  padding: 5px;
  margin: 5px;
  box-sizing: border-box;
  break-inside: avoid;
}

@media (min-width: 320px) and (max-width: 759.5px) {
  .custom_container {
    column-count: 2;
  }
  .video-container {
    width: 320px;
    height: 240px;
  }
}
@media (min-width: 760px) and (max-width: 979.5px) {
  .custom_container {
    column-count: 3;
  }
  .video-container {
    width: 720px;
    height: 540px;
  }
}
@media (min-width: 980px) and (max-width: 1023.5px) {
  .custom_container {
    column-count: 4;
  }
  .video-container {
    width: 760px;
    height: 570px;
  }
}
@media (min-width: 1024px) {
  .custom_container {
    column-count: 5;
  }
  .video-container {
    width: 960px;
    height: 720px;
  }
}
textarea:disabled {
  opacity: 1;
  -webkit-text-fill-color: #000;
}

</style>
<template>
  <div>
    <a-card>
      <a-form-item label="管理模式">
        <a-switch v-model:checked="manageMode"/>
      </a-form-item>
      <a-form :model="searchValue" layout="inline" autocomplete="off">
        <a-form-item label="名称" name="name">
          <a-input v-model:value="searchValue.name"/>
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
          <a-table :dataSource="collections"
                   :row-selection="{ selectedRowKeys: selectedRowKeys, onChange: onSelectChange }" rowKey="id" :columns="columns" :pagination="false">
            <template #bodyCell="{ column, text, record }">
              <template v-if="column.dataIndex === 'thumbnail'">
                <a-image :src="record.cover"
                         :alt="record.name"
                         :width="100"
                />
              </template>
              <template v-if="column.dataIndex === 'actress1'">
                <a-tag v-for="(elem) in parseJson(record.actress)" style="margin-top: 2px" color="green">
                  {{ elem }}
                </a-tag>
              </template>
              <template v-if="column.dataIndex === 'labels1'">
                <a-tag v-for="(elem) in parseJson(record.labels)" style="margin-top: 2px" color="pink">
                  {{ elem }}
                </a-tag>
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
              <template v-for="(collection, index) in collections" :key="index">
                <li class="custom_card" :style="'height:' + config.collectionHeight">
                  <a-image :src="collection.cover"/>
                  <a-row style="margin-top: 5px"><p style="text-overflow: ellipsis; overflow: hidden; margin-bottom: 0">{{ collection.name }}</p></a-row>
                  <a-row style="margin-top: 5px">
                    演员：
                    <a-tag v-for="(elem) in parseJson(collection.actress)" style="margin-top: 2px" color="green">
                      {{ elem }}
                    </a-tag>
                  </a-row>
                  <a-row style="margin-top: 5px">
                    标签：
                    <a-tag v-for="(elem) in parseJson(collection.labels)" style="margin-top: 2px" color="pink">
                      {{ elem }}
                    </a-tag>
                  </a-row>
                  <a-row style="margin-top: 5px">
                    <a-button-group>
                      <a-button size="small" type="primary" @click="updateRecord(collection)">
                        <template #icon>
                          <EditOutlined/>
                        </template>
                      </a-button>
                      <a-popconfirm title="确认删除" @confirm="deleteRecord(collection.id)">
                        <a-button size="small" type="danger">
                          <template #icon>
                            <DeleteOutlined/>
                          </template>
                        </a-button>
                      </a-popconfirm>
                    </a-button-group>
                  </a-row>
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
      <a-form :model="collection" :label-col="{ style: { width: '150px' } }" :wrapper-col="{ span: 14 }">
        <a-form-item label="封面" :rules="[{ required: true }]">
          <a-image v-if="collection.cover" :src="collection.cover" alt="avatar" />
          <a-upload :show-upload-list="false" @change="coverChange" :before-upload="() => false">
            <a-button style="margin: 10px 0 0 10px">
              <plus-outlined/>
            </a-button>
          </a-upload>
        </a-form-item>
        <a-form-item label="名称" :rules="[{ required: true }]">
          <a-input v-model:value="collection.name"/>
        </a-form-item>
        <a-form-item label="演员">
          <a-select
              v-model:value="collection.actress"
              mode="multiple"
              style="width: 100%"
              :filter-option="filterOption"
              :options="actressOptions"
          ></a-select>
        </a-form-item>
        <a-form-item label="标签">
          <a-select
              v-model:value="collection.tag"
              mode="multiple"
              style="width: 100%"
              :filter-option="filterOption"
              :options="tagOptions"
          ></a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, ref} from 'vue';
import {listCollection, createCollection, updateCollection, deleteCollection, detailCollection} from "@/api/collection";
import {DeleteOutlined, EditOutlined, PlusOutlined} from '@ant-design/icons-vue';
import {message} from 'ant-design-vue';
import {optionsActress} from "@/api/actress";
import {optionsTag} from "@/api/tag";
import {parseJson} from "@/utils";
import {getConfig} from "@/api/config";

const manageMode = ref<boolean>(false);
const config = ref<any>({});

const visible = ref<boolean>(false);
const modalTitle = ref<string>('');

const empty = {
  id: 0,
  name: '',
  cover: '',

  actress: [],
  tag: [],
}
const loading = ref<boolean>(false);

const collection = ref<any>(empty);

const coverChange = (info: any) => {
  const { file } = info;
  const reader = new FileReader();
  reader.readAsDataURL(file);
  reader.onload = () => {
    collection.value.cover = reader.result;
  };
};

const create = () => {
  visible.value = true;
  modalTitle.value = '添加对象';
  collection.value = {...empty};
};

const submit = () => {
  let valid = true;
  if (!collection.value.cover) {
    message.warning('请上传封面');
    valid = false;
  }
  if (!collection.value.name) {
    message.warning('请输入集合名称');
    valid = false;
  }

  if (valid) {
    if (collection.value.id === 0) {
      createCollection({
        name: collection.value.name,
        actress: JSON.stringify(collection.value.actress),
        labels: JSON.stringify(collection.value.tag),
        cover: collection.value.cover,
      }).then(() => {
        refresh();
      });
    } else {
      updateCollection({
        id: collection.value.id,
        name: collection.value.name,
        actress: JSON.stringify(collection.value.actress),
        labels: JSON.stringify(collection.value.tag),
        cover: collection.value.cover,
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
  collection.value = {
    id: record.id,
    name: record.name,
    actress: parseJson(record.actress),
    tag: parseJson(record.labels),
    cover: record.cover,
  };
};

const deleteRecord = (id: number) => {
  deleteCollection({
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
  deleteCollection({
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
    title: '名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '演员',
    dataIndex: 'actress1',
  },
  {
    title: '标签',
    dataIndex: 'labels1',
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

const collections = ref<any[]>([]);
const actressOptions = ref<any[]>([]);
const tagOptions = ref<any[]>([]);

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
});

const refresh = () => {
  const data = {
    page: page.value,
    page_size: pageSize.value,
    name: searchValue.value.name,
  };
  listCollection(data).then((res) => {
    collections.value = res.data.data.data;
    total.value = res.data.data.total;
  });
};


onMounted(() => {
  detailCollection(1).then((res) => {
    console.log(res.data.data);
  });

  getConfig().then((res) => {
    config.value = res.data.data;
  });

  optionsActress().then((res) => {
    actressOptions.value = []
    res.data.data.forEach((item: any) => {
      actressMap[item.label] = item.label;
      actressOptions.value.push({
        label: item.label,
        value: item.label,
      });
    });
  });
  optionsTag().then((res) => {
    tagOptions.value = [];
    res.data.data.forEach((item: any) => {
      tagMap[item.label] = item.label;
      tagOptions.value.push({
        label: item.label,
        value: item.label,
      });
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
}
@media (min-width: 760px) and (max-width: 979.5px) {
  .custom_container {
    column-count: 3;
  }
}
@media (min-width: 980px) and (max-width: 1023.5px) {
  .custom_container {
    column-count: 4;
  }
}
@media (min-width: 1024px) {
  .custom_container {
    column-count: 5;
  }
}
textarea:disabled {
  opacity: 1;
  -webkit-text-fill-color: #000;
}
</style>
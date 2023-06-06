<template>
  <div>
    <a-card>
      <a-form-item label="管理模式">
        <a-switch v-model:checked="manageMode"/>
      </a-form-item>
      <a-form :model="searchValue" layout="inline" autocomplete="off">
        <a-form-item label="文件名" name="filename">
          <a-input v-model:value="searchValue.filename"/>
        </a-form-item>
        <a-form-item label="路径" name="path">
          <a-input v-model:value="searchValue.path"/>
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
        <a-form-item label="带有NFO">
          <a-switch v-model:checked="searchValue.nfo"/>
        </a-form-item>

        <a-form-item>
          <a-button type="primary" @click="() => {refresh()}" >搜索</a-button>
        </a-form-item>
      </a-form>
      <a-divider/>
      <a-form-item label="显示">
        <a-space direction="vertical">
          <a-radio-group v-model:value="showImage" :options="[
            { label: '默认', value: 'thumbnail' },
            { label: 'fanart', value: 'fanart' },
            { label: 'poster', value: 'poster' },
            { label: 'thumb', value: 'thumb' },
          ]" />
        </a-space>
      </a-form-item>
      <a-divider/>
      <a-descriptions v-if="collectionInfo.name" :title="'合集名称：' + collectionInfo.name">
        <a-descriptions-item><a-image :width="150" :src="collectionInfo.cover"/></a-descriptions-item>
      </a-descriptions>
      <a-divider v-if="collectionInfo.name"/>
      <a-button-group>
        <a-button type="primary" @click="create">新建对象</a-button>
        <a-button v-if="manageMode"  type="primary" @click="showCollectionModal">添加到集合</a-button>
        <a-button v-if="manageMode"  type="primary" @click="showAddTagModal">添加标签等信息</a-button>
        <a-popconfirm v-if="manageMode && collection" title="确认移除" @confirm="deleteFromCollection">
          <a-button type="danger">从合集中移除</a-button>
        </a-popconfirm>
        <a-popconfirm v-if="manageMode" title="确认删除" @confirm="deleteSelect">
          <a-button type="danger">删除选中</a-button>
        </a-popconfirm>
        <a-popover v-model:visible="pathVisible" title="扫描路径" trigger="click">
          <template #content>
            <a-form :model="object" :label-col="{ style: { width: '75px' } }" :wrapper-col="{ span: 14 }">
              <a-form-item label="路径">
                <a-input v-model:value="object.scanPath" placeholder="路径"/>
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
            </a-form>
            <a-button style="margin-top: 5px" size="middle" @click="scan">确定</a-button>
          </template>
          <a-button type="primary">路径扫描</a-button>
        </a-popover>
        <a-button @click="() => {logVisible = true}">扫描日志</a-button>
        <a-button @click="randomPlay">随机播放</a-button>

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
              <a-image v-if="record.exist_nfo && showImage === 'thumb'"
                  :src="'data:image/jpg;base64,' + record.thumb"
                  :alt="record.name"
                  :width="100"
              />
              <a-image v-if="record.exist_nfo && showImage === 'poster'"
                  :src="'data:image/jpg;base64,' + record.poster"
                  :alt="record.name"
                  :width="100"
              />
              <a-image v-if="record.exist_nfo && showImage === 'fanart'"
                  :src="'data:image/jpg;base64,' + record.fanart"
                  :alt="record.name"
                  :width="100"
              />
              <a-image v-if="!record.exist_nfo || showImage === 'thumbnail'"
                  :src="'data:image/jpg;base64,' + record.thumbnail"
                  :alt="record.name"
                  :width="100"
              />
            </template>
            <template v-else-if="column.dataIndex === 'description'">
              <a href="#" :title="record.description">{{record.description.substring(0, 20)}}</a>
            </template>
            <template v-else-if="column.dataIndex === 'action'">
              <a-button size="small"  @click="playInBrowser(record)">
                <template #icon>
                  <PlayCircleOutlined />
                </template>
              </a-button>
              <a-button size="small"  @click="playInOS(record)">
                <template #icon>
                  <PlaySquareOutlined/>
                </template>
              </a-button>
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
                <li class="custom_card" :style="'height:' + config.objectHeight">
                  <a-image v-if="object.exist_nfo && showImage === 'thumb'" :src="'data:image/jpg;base64,' + object.thumb"/>
                  <a-image v-if="object.exist_nfo && showImage === 'poster'" :src="'data:image/jpg;base64,' + object.poster"/>
                  <a-image v-if="object.exist_nfo && showImage === 'fanart'" :src="'data:image/jpg;base64,' + object.fanart"/>
                  <a-image v-if="!object.exist_nfo || showImage === 'thumbnail'" :src="'data:image/jpg;base64,' + object.thumbnail"/>

                  <a-row style="margin-top: 5px">文件名：<p style="text-overflow: ellipsis; overflow: hidden; margin-bottom: 0">{{ object.name }}</p></a-row>
                  <a-row style="margin-top: 5px">路径：<p style="text-overflow: ellipsis; overflow: hidden; margin-bottom: 0">{{ object.path }}</p></a-row>
                  <a-row style="margin-top: 5px">
                    观看次数：
                    <a-tag color="cyan">
                      {{ object.view_count }}
                    </a-tag>
                  </a-row>
                  <a-row style="margin-top: 5px">
                    演员：
                    <a-tag v-for="(elem) in object.actress" style="margin-top: 2px" color="green">
                      {{ actressMap[elem] }}
                    </a-tag>
                  </a-row>
                  <a-row style="margin-top: 5px">
                    标签：
                    <a-tag v-for="(elem) in object.tag" style="margin-top: 2px" color="pink">
                      {{ tagMap[elem] }}
                    </a-tag>
                  </a-row>
                  <a-row style="margin-top: 5px">
                    分类：
                    <a-tag v-for="(elem) in object.tree" style="margin-top: 2px" color="purple">
                      {{ treeMap[elem] }}
                    </a-tag>
                  </a-row>
                  <a-row style="margin-top: 5px">
                    <a-button-group>
                      <a-button size="small"  @click="playInBrowser(object)">
                        <template #icon>
                          <PlayCircleOutlined />
                        </template>
                      </a-button>
                      <a-button size="small"  @click="playInOS(object)">
                        <template #icon>
                          <PlaySquareOutlined/>
                        </template>
                      </a-button>
                      <a-button size="small" type="primary" @click="updateRecord(object)">
                        <template #icon>
                          <EditOutlined/>
                        </template>
                      </a-button>
                      <a-popconfirm title="确认删除" @confirm="deleteRecord(object.id)">
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
        <a-form-item label="评分">
          <a-input-number id="inputNumber" v-model:value="object.rating" :min="0" :max="10" />
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
    <a-modal v-model:visible="logVisible" title="扫描日志" @ok="" width="1000px" :maskClosable="false" :keyboard="false">
      <template #footer>
        <a-button @click="() => {logVisible = false}">关闭</a-button>
      </template>
      <a-textarea v-model:value="logValue" placeholder="扫描日志" disabled="" :rows="30" />
    </a-modal>

    <a-modal v-model:visible="videoVisible" title="播放" @ok="" @cancel="stopPlay" width="1000px" :maskClosable="false" :keyboard="false">
      <template #footer>
        <a-button @click="stopPlay">关闭</a-button>
      </template>
      <div class="video-container">
        <vue3VideoPlay ref="player" v-bind="options"/>
      </div>
    </a-modal>
    <a-modal v-model:visible="addTagsVisible" title="添加标记，注意仅仅是添加" @ok="addTagsOperation" :maskClosable="false" :keyboard="false">
      <a-form :model="addTagsObject" :label-col="{ style: { width: '150px' } }" :wrapper-col="{ span: 14 }">
        <a-form-item label="演员">
          <a-select
              v-model:value="addTagsObject.actress"
              mode="multiple"
              style="width: 100%"
              :filter-option="filterOption"
              :options="actressOptions"
          ></a-select>
        </a-form-item>
        <a-form-item label="标签">
          <a-select
              v-model:value="addTagsObject.tag"
              mode="multiple"
              style="width: 100%"
              :filter-option="filterOption"
              :options="tagOptions"
          ></a-select>
        </a-form-item>
        <a-form-item label="分组">
          <a-tree-select
              v-model:value="addTagsObject.tree"
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
      </a-form>
    </a-modal>
    <a-modal v-model:visible="add2CollectionVisible" title="添加到集合" @ok="add2Collection" :maskClosable="false" :keyboard="false">
      <a-form :label-col="{ style: { width: '150px' } }" :wrapper-col="{ span: 14 }">
        <a-form-item label="集合">
          <a-select
              v-model:value="collectionId"
              style="width: 100%"
              :filter-option="filterOption"
              :options="collectionOptions"
          ></a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, ref, reactive} from 'vue';
import {createObject, deleteObject, listObject, log, playPath, scanObject, updateObject, videoPath, viewinc, randomPath, addTags} from "@/api/object";
import {DeleteOutlined, EditOutlined, PlaySquareOutlined, PlayCircleOutlined} from '@ant-design/icons-vue';
import {message, TreeSelect} from 'ant-design-vue';
import {optionsActress} from "@/api/actress";
import {optionsTag} from "@/api/tag";
import {optionsTree} from "@/api/tree";
import {buildTree} from "@/utils/util";
import {getConfig} from "@/api/config";
import {optionsCollection, associateCollection, disassociateCollection, getCollection} from "@/api/collection";
import { useRoute } from 'vue-router'

const route = useRoute()

let collection = null;
const collectionInfo = ref<any>({});

const config = ref<any>({});
const showImage = ref<string>("thumbnail");

const videoVisible = ref<boolean>(false);
const player = ref<any>();
const add2CollectionVisible = ref<boolean>(false);
const addTagsVisible = ref<boolean>(false);
const collectionId = ref<any>('');

const options = reactive({
  width: '100%', //播放器高度
  height: '100%', //播放器高度
  color: "#409eff", //主题色
  title: '', //视频名称
  src: "", //视频源
  muted: false, //静音
  webFullScreen: false,
  poster: "", //封面图
  speedRate: ["0.75", "1.0", "1.25", "1.5", "2.0"], //播放倍速
  autoPlay: false, //自动播放
  loop: false, //循环播放
  mirror: false, //镜像画面
  ligthOff: false,  //关灯模式
  volume: 0.2, //默认音量大小
  control: true, //是否显示控制
  controlBtns:['audioTrack', 'quality', 'speedRate', 'volume', 'setting', 'pip', 'pageFullScreen', 'fullScreen'] //显示所有按钮,
})

const SHOW_ALL = TreeSelect.SHOW_ALL;

const pathVisible = ref<boolean>(false);
const logVisible = ref<boolean>(false);

const logValue = ref<string>('');

const scanLog = () => {
  setTimeout(() => {
    if (logVisible.value) {
      log().then(response => {
        logValue.value += response.data.data;
      });
    }
    scanLog()
  }, 3000);
};
scanLog();


const scan = () => {
  if (!object.value.scanPath) {
    message.warning('请输入扫描路径');
    return;
  }
  pathVisible.value = false;
  scanObject({
    path: object.value.scanPath,
    actress: object.value.actress,
    tag: object.value.tag,
    tree: object.value.tree,
  }).then(response => {
    message.success('任务开始');
    logVisible.value = true;
  });
};

const manageMode = ref<boolean>(false);

const visible = ref<boolean>(false);
const modalTitle = ref<string>('');

const empty = {
  id: 0,
  name: '',
  rating: 0,
  description: '',
  path: '',
  magnet: '',
  scanPath: '',

  actress: [],
  tag: [],
  tree: [],
}

const object = ref<any>(empty);

const addTagsObject = ref<any>(empty);

const create = () => {
  visible.value = true;
  modalTitle.value = '添加对象';
  object.value = {...empty};
  addTagsObject.value = {...empty};
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
        rating: object.value.rating,
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
        rating: object.value.rating,
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

const playInOS = (record: any) => {
  playPath(record.path).then(response => {
    message.info("播放成功");
  });
  viewinc(record.id).then(response => {});
};

const playInBrowser = (record: any) => {
  videoVisible.value = true;
  options.title = record.name;
  options.poster = 'data:image/jpg;base64,' + record.thumbnail;
  options.src = videoPath(record.path);
  viewinc(record.id).then(response => {});
};

const randomPlay = () => {
  randomPath().then(response => {
    options.title = response.data.data.name;
    options.poster = 'data:image/jpg;base64,' + response.data.data.thumbnail;
    options.src = videoPath(response.data.data.path);
    videoVisible.value = true;
  });
};

const stopPlay = () => {
  videoVisible.value = false;
  player.value.pause()
};

const updateRecord = (record: any) => {
  visible.value = true;
  modalTitle.value = '修改对象';
  object.value = {
    id: record.id,
    name: record.name,
    rating: record.rating,
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

const deleteFromCollection = () => {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择要删除的对象');
    return;
  }
  disassociateCollection({
    collection_id: collection,
    object_id: selectedRowKeys.value
  }).then(() => {
    refresh();
  });
};

const showCollectionModal = () => {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择要添加的对象');
    return;
  }
  add2CollectionVisible.value = true
}

const showAddTagModal = () => {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择要添加的对象');
    return;
  }
  addTagsVisible.value = true
}

const add2Collection = () => {
  if (collectionId.value === '') {
    message.warning('请选择要添加的合集');
    return;
  }

  associateCollection({
    collection_id: collectionId.value,
    object_id: selectedRowKeys.value
  }).then(() => {
    message.success("添加成功");
    add2CollectionVisible.value = false;
  });
}

const addTagsOperation = () => {
  addTags({
    object_id: selectedRowKeys.value,
    tags: addTagsObject.value.tag,
    trees: addTagsObject.value.tree,
    actresses: addTagsObject.value.actress,
  }).then(() => {
    message.success("添加成功");
    addTagsVisible.value = false;
    addTagsObject.value = {...empty};
    refresh();
  });
}

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
    title: '评分',
    dataIndex: 'rating',
    key: 'rating',
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
const collectionOptions = ref<any[]>([]);

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
  filename: '',
  path: '',
  actress: [],
  tag: [],
  tree: [],
  nfo: false,
});

const refresh = () => {
  const data = {
    page: page.value,
    page_size: pageSize.value,
    filename: searchValue.value.filename,
    path: searchValue.value.path,
    actress: searchValue.value.actress,
    tag: searchValue.value.tag,
    tree: searchValue.value.tree,
    nfo: searchValue.value.nfo,
    collection: route.query.collection,
  };
  listObject(data).then((res) => {
    objects.value = res.data.data.data;
    total.value = res.data.data.total;
  });
};

onMounted(() => {
  if (route.query.collection) {
    collection = route.query.collection;
    getCollection(collection).then((res) => {
      collectionInfo.value = res.data.data
    });
  }

  getConfig().then((res) => {
    config.value = res.data.data;
  });
  optionsCollection().then((res) => {
    collectionOptions.value = res.data.data;
  });
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
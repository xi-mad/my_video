export const buildTree = (data: any[]) => {
    const tree: any[] = []
    const treeMap = {};
    data.forEach((item) => {
        treeMap[item.value] = item;
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
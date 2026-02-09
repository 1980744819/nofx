// 声明式 Pipeline
pipeline {
    // 选择执行节点（若有多个 agent，可指定标签；默认 master 节点）
    agent {
        kubernetes {
            cloud 'itx'  // 对应 Jenkins 后台配置的 K8s 集群名称（不变）
            inheritFrom 'jenkins-agent'  // 核心：引用后台创建的静态 Pod 模板名称（替代废弃的 label）
        }
    }
    
    // 环境变量定义（统一配置，便于修改）
    environment {
        // 1. 代码仓库配置
        GIT_REPO_URL = "https://github.com/1980744819/nofx.git"
        GIT_BRANCH = "main"  // 要拉取的分支（如 main、dev）
        
        // 2. 镜像配置
        DOCKERFILE_PATH = "./docker"  // Dockerfile 所在路径（根目录为 "."，子目录如 "./docker"）
        REPO_ADDR = "192.168.1.7:30002"  // 私人镜像仓库地址+镜像名（替换为你的实际地址）
        PROJECT = "nofx"
        FRONT_IMAGE_REPO = "${REPO_ADDR}/${PROJECT}/frontend"  // 私人镜像仓库地址+镜像名    
        BACK_IMAGE_REPO = "${REPO_ADDR}/${PROJECT}/backend"  // 私人镜像仓库地址+镜像名    
        NAMESPACE = "prod"
        GIT_CREDENTIAL_ID = "credential_github"  // GitHub 凭证 ID
        HARBOR_CREDENTIAL_ID = "credential_harbor_admin"  // Harbor 凭证 ID
    }
    
    // 流水线阶段定义（拉取代码 → 构建镜像 → 推送镜像）
    stages {
        // 阶段 1：拉取 GitHub 代码仓库
        stage('Pull GitHub Code') {
            steps {
                echo "开始拉取 GitHub 仓库 ${GIT_REPO_URL} 分支 ${GIT_BRANCH}..."
                // 使用 git 步骤拉取代码，关联已配置的 GitHub 凭证
                git(
                    url: "${GIT_REPO_URL}",
                    branch: "${GIT_BRANCH}",
                    credentialsId: "${GIT_CREDENTIAL_ID}",  // GitHub 凭证 ID
                    changelog: false,
                    poll: false
                )
                script {
                    if (!env.GIT_COMMIT || env.GIT_COMMIT == 'null') {
                        env.GIT_COMMIT = sh(returnStdout: true, script: "git rev-parse HEAD").trim()
                    }
                    env.GIT_SHORT = sh(returnStdout: true, script: "git rev-parse --short HEAD").trim()
                    echo "Commit: ${env.GIT_COMMIT}, Short: ${env.GIT_SHORT}"
                }
                echo "代码拉取完成！"
            }
        }
        stage('Init Image Variables') {
            steps {
                script {
                    // 此时已拉取代码，env.GIT_COMMIT 可用，顺序赋值无冲突
                    IMAGE_TAG = "${env.GIT_COMMIT.substring(0, 8)}"
                    FRONT_FULL_IMAGE_NAME = "${FRONT_IMAGE_REPO}:${IMAGE_TAG}"
                    BACK_FULL_IMAGE_NAME = "${BACK_IMAGE_REPO}:${IMAGE_TAG}"
                    
                    // （可选）将变量注入全局环境，后续阶段可直接使用
                    env.IMAGE_TAG = IMAGE_TAG
                    env.FRONT_FULL_IMAGE_NAME = FRONT_FULL_IMAGE_NAME
                    env.BACK_FULL_IMAGE_NAME = BACK_FULL_IMAGE_NAME
                    
                    echo "初始化镜像变量完成："
                    echo "镜像标签：${IMAGE_TAG}"
                    echo "前端完整镜像名：${FRONT_FULL_IMAGE_NAME}"
                    echo "后端完整镜像名：${BACK_FULL_IMAGE_NAME}"
                }
            }
        }
        // 阶段 2：构建 Docker 镜像
        stage('Build Frontend podman Image') {
            steps {
                container('podman') {
                    withCredentials([usernamePassword(credentialsId: "${HARBOR_CREDENTIAL_ID}", usernameVariable: 'HARBOR_USER', passwordVariable: 'HARBOR_PASS')]) {
                        sh '''
                            set -e
                            
                            # 登录 Harbor 镜像仓库
                            echo "登录 Harbor 镜像仓库..."
                            podman login --tls-verify=false -u "${HARBOR_USER}" -p "${HARBOR_PASS}" "${REPO_ADDR}"
                        '''
                    }

                    echo "开始构建前端镜像 ${FRONT_FULL_IMAGE_NAME}..."
                    // 使用 podman build 构建镜像，指定 Dockerfile 路径和镜像标签
                    sh """
                        set -e
                        podman build -f ${DOCKERFILE_PATH}/Dockerfile.frontend \
                        --squash --network=host \
                        --tls-verify=false \
                        -t ${FRONT_FULL_IMAGE_NAME} \
                        .
                    """
                    echo "前端镜像构建完成！"
                }
            }
        }
        stage('Build Backend podman Image') {
            steps {
                container('podman') {
                    echo "开始构建后端镜像 ${BACK_FULL_IMAGE_NAME}..."
                    // 使用 podman build 构建镜像，指定 Dockerfile 路径和镜像标签
                    sh """
                        set -e
                        podman build -f ${DOCKERFILE_PATH}/Dockerfile.backend \
                        --squash --network=host \
                        --tls-verify=false \
                        -t ${BACK_FULL_IMAGE_NAME} \
                        .
                    """
                    echo "后端镜像构建完成！"
                }
            }
        }
        
        // 阶段 3：登录私人镜像仓库并推送镜像
        stage('Push to Private podman Repo') {
            steps {
                container('podman') {
                    withCredentials([usernamePassword(credentialsId: "${HARBOR_CREDENTIAL_ID}", usernameVariable: 'HARBOR_USER', passwordVariable: 'HARBOR_PASS')]) {
                        sh '''
                            set -e
                            
                            # 登录 Harbor
                            echo "登录 Harbor 镜像仓库..."
                            podman login --tls-verify=false -u "${HARBOR_USER}" -p "${HARBOR_PASS}" "${REPO_ADDR}"
                            
                            # 推送镜像
                            echo "推送前端镜像 ${FRONT_FULL_IMAGE_NAME} 到 Harbor..."
                            podman push --tls-verify=false "${FRONT_FULL_IMAGE_NAME}"
                            echo "前端镜像推送完成！"
                            
                            echo "推送后端镜像 ${BACK_FULL_IMAGE_NAME} 到 Harbor..."
                            podman push --tls-verify=false "${BACK_FULL_IMAGE_NAME}"
                            echo "后端镜像推送完成！"
                        '''
                    }
                }
            }
        }
        stage('helm deploy') {
            steps {
                container('helm') {
                    sh '''
                    set -euo pipefail

                    helm upgrade --install nofx ./chart \
                        --namespace ${NAMESPACE} --create-namespace \
                        --set backend.image.repository=${BACK_IMAGE_REPO} \
                        --set backend.image.tag=${IMAGE_TAG} \
                        --set frontend.image.repository=${FRONT_IMAGE_REPO} \
                        --set frontend.image.tag=${IMAGE_TAG}
                    '''
                }
            }
        }
    }
    
    // 后置操作（无论成功/失败都执行，可选）
    post {
        success {
            script {
                container('podman') {
                    echo "流水线构建成功，已推送镜像到私人仓库！"
                    echo "前端镜像：${FRONT_FULL_IMAGE_NAME}"
                    echo "后端镜像：${BACK_FULL_IMAGE_NAME}"
                    // 清理本地构建的镜像
                    echo "开始清理本地构建的镜像..."
                    sh """
                        podman rmi ${FRONT_FULL_IMAGE_NAME} || true
                        podman rmi ${BACK_FULL_IMAGE_NAME} || true
                    """
                    echo "本地镜像清理完成！"
                }
            }
        }
        failure {
            echo "流水线构建失败，请查看日志排查问题！"
        }
    }
}
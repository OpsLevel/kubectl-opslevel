package common

import (
	"github.com/opslevel/opslevel-go/v2023"
	"github.com/rocktavious/autopilot/v2023"
	"testing"
)

var k8sResource = `{
    "apiVersion": "apps/v1",
    "kind": "Deployment",
    "metadata": {
        "annotations": {
            "deployment.kubernetes.io/revision": "243",
            "kots.io/app-slug": "opslevel",
			"opslevel.com/description": "this is a description",
			"opslevel.com/owner": "velero",
			"opslevel.com/lifecycle": "alpha",
			"opslevel.com/tier": "tier_1",
			"opslevel.com/product": "jklabs",
			"opslevel.com/language": "ruby",
			"opslevel.com/framework": "rails",
			"opslevel.com/system": "monolith",
			"opslevel.com/tools.logs.my-logs": "https://example.com",
			"opslevel.com/repo.terraform.clusters.dev.opslevel": "gitlab.com:opslevel/terraform",
            "repo": "github.com:hashicorp/vault"
        },
        "creationTimestamp": "2023-07-19T18:04:03Z",
        "generation": 243,
        "labels": {
            "app.kubernetes.io/instance": "web",
            "app.kubernetes.io/part-of": "opslevel",
            "kots.io/app-slug": "opslevel",
            "kots.io/backup": "velero"
        },
        "name": "web",
        "namespace": "self-hosted",
        "resourceVersion": "383293724",
        "uid": "19d729de-f708-437c-8b65-10fa06d5dfd5"
    },
    "spec": {
        "progressDeadlineSeconds": 600,
        "replicas": 2,
        "revisionHistoryLimit": 3,
        "selector": {
            "matchLabels": {
                "app.kubernetes.io/instance": "web",
                "app.kubernetes.io/part-of": "opslevel"
            }
        },
        "strategy": {
            "rollingUpdate": {
                "maxSurge": "1",
                "maxUnavailable": 0
            },
            "type": "RollingUpdate"
        },
        "template": {
            "metadata": {
                "annotations": {
                    "kots.io/app-slug": "opslevel"
                },
                "creationTimestamp": null,
                "labels": {
                    "app.kubernetes.io/instance": "web",
                    "app.kubernetes.io/part-of": "opslevel",
					"environment": "dev",
                    "collect-logs": "true"
                }
            },
            "spec": {
                "containers": [
                    {
                        "args": [
                            "bundle",
                            "exec",
                            "puma",
                            "-C ./config/puma.rb"
                        ],
                        "env": [],
                        "envFrom": [
                            {
                                "configMapRef": {
                                    "name": "opslevel"
                                }
                            },
                            {
                                "secretRef": {
                                    "name": "opslevel"
                                }
                            }
                        ],
                        "image": "opslevel/opslevel:main-240131e5",
                        "imagePullPolicy": "Always",
                        "lifecycle": {
                            "preStop": {
                                "exec": {
                                    "command": [
                                        "sleep",
                                        "15"
                                    ]
                                }
                            }
                        },
                        "livenessProbe": {
                            "failureThreshold": 3,
                            "initialDelaySeconds": 3,
                            "periodSeconds": 20,
                            "successThreshold": 1,
                            "tcpSocket": {
                                "port": "opslevel"
                            },
                            "timeoutSeconds": 1
                        },
                        "name": "web",
                        "ports": [
                            {
                                "containerPort": 3000,
                                "name": "opslevel",
                                "protocol": "TCP"
                            }
                        ],
                        "readinessProbe": {
                            "failureThreshold": 3,
                            "httpGet": {
                                "path": "/api/ping",
                                "port": "opslevel",
                                "scheme": "HTTP"
                            },
                            "initialDelaySeconds": 5,
                            "periodSeconds": 10,
                            "successThreshold": 2,
                            "timeoutSeconds": 1
                        },
                        "resources": {
                            "limits": {
                                "cpu": "1",
                                "memory": "1536Mi"
                            },
                            "requests": {
                                "cpu": "500m",
                                "memory": "500Mi"
                            }
                        },
                        "terminationMessagePath": "/dev/termination-log",
                        "terminationMessagePolicy": "File"
                    }
                ],
                "dnsPolicy": "ClusterFirst",
                "imagePullSecrets": [
                    {
                        "name": "opslevel-registry"
                    }
                ],
                "initContainers": [
                    {
                        "args": [
                            "bundle",
                            "exec",
                            "rake",
                            "db:abort_if_pending_migrations"
                        ],
                        "envFrom": [
                            {
                                "configMapRef": {
                                    "name": "opslevel"
                                }
                            },
                            {
                                "secretRef": {
                                    "name": "opslevel"
                                }
                            }
                        ],
                        "image": "opslevel/opslevel:main-240131e5",
                        "imagePullPolicy": "Always",
                        "name": "migrations",
                        "resources": {},
                        "terminationMessagePath": "/dev/termination-log",
                        "terminationMessagePolicy": "File"
                    }
                ],
                "nodeSelector": {
                    "kubernetes.io/os": "linux"
                },
                "restartPolicy": "Always",
                "schedulerName": "default-scheduler",
                "securityContext": {},
                "terminationGracePeriodSeconds": 30,
                "topologySpreadConstraints": [
                    {
                        "labelSelector": {
                            "matchLabels": {
                                "app.kubernetes.io/name": "web",
                                "app.kubernetes.io/part-of": "opslevel"
                            }
                        },
                        "maxSkew": 1,
                        "topologyKey": "topology.kubernetes.io/zone",
                        "whenUnsatisfiable": "ScheduleAnyway"
                    }
                ]
            }
        }
    },
    "status": {
        "availableReplicas": 2,
        "conditions": [
            {
                "lastTransitionTime": "2023-07-19T18:05:54Z",
                "lastUpdateTime": "2023-07-19T18:05:54Z",
                "message": "Deployment has minimum availability.",
                "reason": "MinimumReplicasAvailable",
                "status": "True",
                "type": "Available"
            },
            {
                "lastTransitionTime": "2023-08-31T09:00:52Z",
                "lastUpdateTime": "2023-09-25T16:26:30Z",
                "message": "ReplicaSet \"web-6fd48cb855\" has successfully progressed.",
                "reason": "NewReplicaSetAvailable",
                "status": "True",
                "type": "Progressing"
            }
        ],
        "observedGeneration": 243,
        "readyReplicas": 2,
        "replicas": 2,
        "updatedReplicas": 2
    }
}
`

func TestJQServicePArserSimpleConfig(t *testing.T) {
	// Arrange
	config, err := GetConfig(ConfigSimple)
	if err != nil {
		t.Error(err)
	}
	parser := NewJQServiceParser(config.Service.Import[0].OpslevelConfig)
	// Act
	service, err := parser.Run(k8sResource)
	if err != nil {
		t.Error(err)
	}
	// Assert
	autopilot.Equals(t, "web", service.Name)
	autopilot.Equals(t, "self-hosted", service.Owner)
	autopilot.Equals(t, "", service.Lifecycle)
	autopilot.Equals(t, "", service.Tier)
	autopilot.Equals(t, "", service.Product)
	autopilot.Equals(t, "", service.Language)
	autopilot.Equals(t, "", service.Framework)
	autopilot.Equals(t, "", service.System)
	autopilot.Equals(t, 1, len(service.Aliases))
	autopilot.Equals(t, "k8s:web-self-hosted", service.Aliases[0])
	autopilot.Equals(t, 1, len(service.TagCreates))
	autopilot.Equals(t, opslevel.TagInput{Key: "environment", Value: "dev"}, service.TagCreates[0])
	autopilot.Equals(t, 5, len(service.TagAssigns))
	autopilot.Equals(t, opslevel.TagInput{Key: "imported", Value: "kubectl-opslevel"}, service.TagAssigns[0])
	autopilot.Equals(t, 0, len(service.Tools))
	autopilot.Equals(t, 0, len(service.Repositories))
}

func TestJQServiceParserSampleConfig(t *testing.T) {
	// Arrange
	config, err := GetConfig(ConfigSample)
	if err != nil {
		t.Error(err)
	}

	parser := NewJQServiceParser(config.Service.Import[0].OpslevelConfig)
	// Act
	service, err := parser.Run(k8sResource)
	if err != nil {
		t.Error(err)
	}
	// Assert
	autopilot.Equals(t, "web", service.Name)
	autopilot.Equals(t, "this is a description", service.Description)
	autopilot.Equals(t, "velero", service.Owner)
	autopilot.Equals(t, "alpha", service.Lifecycle)
	autopilot.Equals(t, "tier_1", service.Tier)
	autopilot.Equals(t, "jklabs", service.Product)
	autopilot.Equals(t, "ruby", service.Language)
	autopilot.Equals(t, "rails", service.Framework)
	autopilot.Equals(t, "monolith", service.System)
	autopilot.Equals(t, "k8s:web-self-hosted", service.Aliases[0])
	autopilot.Equals(t, "self-hosted-web", service.Aliases[1])
	autopilot.Equals(t, 1, len(service.TagCreates))
	autopilot.Equals(t, opslevel.TagInput{Key: "environment", Value: "dev"}, service.TagCreates[0])
	autopilot.Equals(t, 5, len(service.TagAssigns))
	autopilot.Equals(t, opslevel.TagInput{Key: "imported", Value: "kubectl-opslevel"}, service.TagAssigns[0])
	autopilot.Equals(t, 1, len(service.Tools))
	autopilot.Equals(t, 3, len(service.Repositories))
}

func BenchmarkJQParser_New(b *testing.B) {
	config, _ := GetConfig(ConfigSample)
	parser := NewJQServiceParser(config.Service.Import[0].OpslevelConfig)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.Run(k8sResource)
	}
}

